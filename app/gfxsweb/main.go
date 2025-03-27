package main

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"encoding/base64"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "image/gif"
	_ "image/jpeg"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	gfxi "github.com/toxyl/gfx/image"
	"github.com/toxyl/gfx/parser"
)

//go:embed index.html
var index string

//go:embed style.css
var style string

//go:embed script.js
var script string

//go:embed gfxs-mode.js
var gfxsmode string

var MAX_MP int = 2048 * 1536

// getCurrentImagePath returns the path of the current static image.
func getCurrentImagePath() (string, error) {
	if _, err := os.Stat("image.png"); err == nil {
		return "image.png", nil
	} else if _, err := os.Stat("image.jpg"); err == nil {
		return "image.jpg", nil
	}
	return "", fmt.Errorf("no image file found")
}

// updateDimensionsIfMissing checks the [COMPOSITION] section and, if width/height are missing,
// appends the missing lines using the provided dimensions.
func updateDimensionsIfMissing(gfxs string, width, height int) string {
	str := strings.Split(gfxs, "[COMPOSITION]")
	if len(str) < 2 {
		return gfxs
	}

	pre := str[0] + "\n\n[COMPOSITION]\n"
	post := str[1]

	if ok, _ := regexp.MatchString(`(?m)^\s*width\s*=`, post); !ok {
		pre += "\nwidth = " + strconv.Itoa(width)
	}
	if ok, _ := regexp.MatchString(`(?m)^\s*height\s*=`, post); !ok {
		pre += "\nheight = " + strconv.Itoa(height)
	}

	return pre + post
}

// changeExtension returns the base name with a .png extension.
func changeExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	base := strings.TrimSuffix(filename, ext)
	return base + ".png"
}

func saveBytesAsPNG(filename string, buf []byte, maxMP int) error {
	_, format, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("invalid image file: %s", err.Error())
	}

	if format != "png" && format != "jpeg" {
		return fmt.Errorf("unsupported image format: %s", format)
	}

	i := gfxi.NewFromBytes(format, buf).ResizeToMaxMP(maxMP)
	i.SaveAsPNG(filename)
	return nil
}

var (
	isRendering = false
	rtIndex     = func(c *fiber.Ctx) error {
		c.Context().SetContentType("text/html")
		return c.SendString(index)
	}
	rtStyle = func(c *fiber.Ctx) error {
		c.Context().SetContentType("text/css")
		return c.SendString(style)
	}
	rtScript = func(c *fiber.Ctx) error {
		c.Context().SetContentType("text/javascript")
		return c.SendString(script)
	}
	rtSyntax = func(c *fiber.Ctx) error {
		c.Context().SetContentType("text/javascript")
		return c.SendString(gfxsmode)
	}
	rtOriginal = func(c *fiber.Ctx) error {
		imagePath, err := getCurrentImagePath()
		if err != nil {
			return c.Status(http.StatusNotFound).SendString("No image file found")
		}
		return c.SendFile(imagePath)
	}
	rtProcessed = func(c *fiber.Ctx) error {
		for isRendering {
			time.Sleep(100 * time.Millisecond)
		}
		if _, err := os.Stat("processed.png"); err != nil {
			return c.Status(http.StatusNotFound).SendString("No processed image available")
		}
		return c.SendFile("processed.png")
	}
	rtUpload = func(c *fiber.Ctx) error {
		fileHeader, err := c.FormFile("image")
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString("No file uploaded")
		}
		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString("Failed to open uploaded file: " + err.Error())
		}
		defer file.Close()

		buf, err := io.ReadAll(file)
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString("Failed to read file: " + err.Error())
		}

		filename := "image.png"
		if err := saveBytesAsPNG(filename, buf, MAX_MP); err != nil {
			return c.Status(http.StatusBadRequest).SendString("Invalid image file: " + err.Error())
		}
		return c.SendString("OK")
	}
	rtRender = func(c *fiber.Ctx) error {
		isRendering = true
		defer func() { isRendering = false }()
		gfxs := c.FormValue("gfxs")
		if gfxs == "" {
			return c.Status(http.StatusBadRequest).SendString("GFXS script is required")
		}

		// Get the original image path.
		imagePath, err := getCurrentImagePath()
		if err != nil {
			return c.Status(http.StatusNotFound).SendString("No image file found")
		}

		// Open the image file to get its dimensions.
		file, err := os.Open(imagePath)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Failed to open image file: " + err.Error())
		}
		defer file.Close()
		img, _, err := image.Decode(file)
		if err == nil {
			width := img.Bounds().Dx()
			height := img.Bounds().Dy()
			gfxs = updateDimensionsIfMissing(gfxs, width, height)
		}

		absImagePath, err := filepath.Abs(imagePath)
		if err != nil {
			absImagePath = imagePath
		}
		gfxs = strings.ReplaceAll(gfxs, `$IMG`, absImagePath)

		comp, err := parser.ParseComposition(gfxs)
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString("GFXS script could not be parsed: " + err.Error())
		}

		// Render the composition and encode the processed image as PNG into memory.
		outBuffer := new(bytes.Buffer)
		renderedComp := comp.Render()
		renderedData := renderedComp.Get()
		if err := png.Encode(outBuffer, renderedData); err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Failed to encode processed image: " + err.Error())
		}

		// Encode the processed image to base64.
		processedBase64 := base64.StdEncoding.EncodeToString(outBuffer.Bytes())

		// Read the original image file and encode it.
		origBytes, err := os.ReadFile(imagePath)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Failed to read original image: " + err.Error())
		}
		originalBase64 := base64.StdEncoding.EncodeToString(origBytes)

		// Return JSON containing both base64-encoded images.
		return c.JSON(fiber.Map{
			"update":    true,
			"original":  originalBase64,
			"processed": processedBase64,
		})
	}
	rtRenderBatch = func(c *fiber.Ctx) error {
		isRendering = true
		defer func() { isRendering = false }()
		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString("Failed to parse multipart form: " + err.Error())
		}
		gfxs := c.FormValue("gfxs")
		if gfxs == "" {
			return c.Status(http.StatusBadRequest).SendString("GFXS script is required")
		}
		files := form.File["images"]
		if len(files) == 0 {
			return c.Status(http.StatusBadRequest).SendString("No images provided")
		}

		tempDir, err := os.MkdirTemp("", "batchrender")
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Failed to create temp directory: " + err.Error())
		}
		defer os.RemoveAll(tempDir)

		type processedFile struct {
			originalName  string
			processedPath string
		}
		var processedFiles []processedFile

		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				continue
			}

			buf, err := io.ReadAll(file)
			if err != nil {
				defer file.Close()
				return c.Status(http.StatusBadRequest).SendString("Failed to read file: " + err.Error())
			}
			file.Close()

			tempOrigPath := strings.TrimSuffix(filepath.Join(tempDir, strings.ReplaceAll(fileHeader.Filename, " ", "_")), ".png") + ".png"

			if err := saveBytesAsPNG(tempOrigPath, buf, MAX_MP); err != nil {
				return c.Status(http.StatusBadRequest).SendString("Failed to save image: " + err.Error())
			}

			i := gfxi.NewFromFile(tempOrigPath)
			updatedGfxs := strings.ReplaceAll(updateDimensionsIfMissing(gfxs, i.W(), i.H()), `$IMG`, i.Path())
			comp, err := parser.ParseComposition(updatedGfxs)
			if err != nil {
				continue
			}
			processedName := changeExtension(fileHeader.Filename)
			processedPath := filepath.Join(tempDir, processedName)
			comp.Render().SaveAsPNG(processedPath)
			processedFiles = append(processedFiles, processedFile{originalName: processedName, processedPath: processedPath})
		}

		var bufZip bytes.Buffer
		zipWriter := zip.NewWriter(&bufZip)
		for _, pf := range processedFiles {
			data, err := os.ReadFile(pf.processedPath)
			if err != nil {
				continue
			}
			f, err := zipWriter.Create(pf.originalName)
			if err != nil {
				continue
			}
			_, err = f.Write(data)
			if err != nil {
				continue
			}
		}
		zipWriter.Close()
		c.Response().Header.Set("Content-Type", "application/zip")
		c.Response().Header.Set("Content-Disposition", "attachment; filename=batch.zip")
		return c.Send(bufZip.Bytes())
	}
	rtFilters = func(c *fiber.Ctx) error {
		files, err := os.ReadDir("filters")
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Failed to read filters directory: " + err.Error())
		}
		var names []string
		for _, f := range files {
			if !f.IsDir() {
				names = append(names, strings.TrimSuffix(f.Name(), ".gfxs"))
			}
		}
		return c.JSON(names)
	}
	rtFilter = func(c *fiber.Ctx) error {
		name := c.Query("name")
		if name == "" {
			return c.Status(http.StatusBadRequest).SendString("Filter name required")
		}
		path := filepath.Join("filters", name+".gfxs")
		data, err := os.ReadFile(path)
		if err != nil {
			return c.Status(http.StatusNotFound).SendString("Filter not found")
		}
		return c.SendString(string(data))
	}
	rtSaveFilter = func(c *fiber.Ctx) error {
		name := c.FormValue("name")
		filterContent := c.FormValue("filter")
		if name == "" || filterContent == "" {
			return c.Status(http.StatusBadRequest).SendString("Name and filter content are required")
		}
		path := filepath.Join("filters", name+".gfxs")
		err := os.WriteFile(path, []byte(filterContent), 0644)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Failed to save filter: " + err.Error())
		}
		return c.SendString("OK")
	}
	rtRenderFilter = func(c *fiber.Ctx) error {
		// config
		baseImage := "render.png"

		// load composition
		name := c.Query("name")
		if name == "" {
			return c.Status(http.StatusBadRequest).SendString("Filter name required")
		}
		path := filepath.Join("filters", name+".gfxs")
		data, err := os.ReadFile(path)
		if err != nil {
			return c.Status(http.StatusNotFound).SendString("Filter not found")
		}

		// process url
		url := c.Query("url")
		if url == "" {
			return c.Status(http.StatusBadRequest).SendString("URL required")
		}
		u, err := base64.URLEncoding.DecodeString(url)
		if err != nil || len(u) == 0 {
			return c.Status(http.StatusBadRequest).SendString("Valid URL required")
		}
		i := gfxi.NewFromURL(string(u)).ResizeToMaxMP(MAX_MP).SaveAsPNG(baseImage)
		filterText := updateDimensionsIfMissing(strings.ReplaceAll(string(data), `$IMG`, i.Path()), i.W(), i.H())
		comp, err := parser.ParseComposition(filterText)
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString("GFXS filter could not be parsed: " + err.Error())
		}
		_ = comp.Render().SaveAsPNG(baseImage)
		time.Sleep(100 * time.Millisecond)
		return c.SendFile(baseImage)
	}
)

func main() {
	fport := flag.Uint("p", 8080, "The port to run the server on, defaults to 8080.")
	flag.Parse()

	os.MkdirAll("filters", 0755)

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024 * 1024,
	})
	app.Use(recover.New())
	app.Use(cors.New())

	app.Get("/", rtIndex)
	app.Get("/styles.css", rtStyle)
	app.Get("/script.js", rtScript)
	app.Get("/gfxs-mode.js", rtSyntax)
	app.Get("/original", rtOriginal)
	app.Get("/processed", rtProcessed)
	app.Post("/upload", rtUpload)
	app.Post("/render", rtRender)
	app.Post("/renderBatch", rtRenderBatch)
	app.Get("/filters", rtFilters)
	app.Get("/filter", rtFilter)
	app.Post("/saveFilter", rtSaveFilter)
	app.Get("/renderFilter", rtRenderFilter)

	app.Listen(":" + fmt.Sprint(*fport))
}
