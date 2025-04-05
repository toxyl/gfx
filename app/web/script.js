document.addEventListener("DOMContentLoaded", () => {
  // Initialize CodeMirror editors.
  const editorFilters = CodeMirror.fromTextArea(document.getElementById("gfxse-filters"), {
    lineNumbers: true,
    mode: "gfxscript",
  });
  const editorComposition = CodeMirror.fromTextArea(document.getElementById("gfxse-composition"), {
    lineNumbers: true,
    mode: "gfxscript",
  });
  const editorVars = CodeMirror.fromTextArea(document.getElementById("gfxse-vars"), {
    lineNumbers: true,
    mode: "gfxscript",
  });
  const editorLayers = CodeMirror.fromTextArea(document.getElementById("gfxse-layers"), {
    lineNumbers: true,
    mode: "gfxscript",
  });

  // Filter management controls.
  const filterDropdown = document.getElementById("filter-dropdown");
  const saveFilterBtn = document.getElementById("save-filter-btn");
  const forceRenderBtn = document.getElementById("force-render-btn");

  // Preview mode controls and elements.
  const tabButtons = document.querySelectorAll("#preview-tabs .tab-button");
  const previewSideBySide = document.getElementById("preview-side-by-side");
  const previewBoth = document.getElementById("preview-both");
  const previewImage = document.getElementById("preview-image"); // For "both" mode.
  const processedImageSB = document.getElementById("processed-image"); // Side-by-side: right pane.
  const originalImageSB = document.getElementById("original-image");   // Side-by-side: left pane.
  const autoRenderTimerDisplay = document.getElementById("auto-render-timer");

  // Overlays for batch processing.
  const overlayFull = document.getElementById("overlay-full");

  let autoRenderTimeoutId = null;
  let autoRenderScheduledTime = null;
  let isRenderInProgress = false;
  let lastRenderedContent = "";
  let currentTimestamp = Date.now();
  const AUTO_RENDER_DELAY = 10000; // 10 seconds

  // Global variables to store the latest base64 image data.
  let base64Original = "";
  let base64Processed = "";

  previewImage.addEventListener("mousedown", () => {
    previewImage.src = "data:image/png;base64," + base64Original;
  });
  previewImage.addEventListener("mouseup", () => {
    previewImage.src = "data:image/png;base64," + base64Processed;
  });
  previewImage.addEventListener("mouseleave", () => {
    previewImage.src = "data:image/png;base64," + base64Processed;
  });

  // --- Helper: Show and hide overlays ---
  function showOverlay(overlayElement) {
    overlayElement.style.display = "flex";
  }
  function hideOverlay(overlayElement) {
    overlayElement.style.display = "none";
  }

  // --- Tab Switching for Preview Mode ---
  tabButtons.forEach(button => {
    button.addEventListener("click", () => {
      tabButtons.forEach(btn => btn.classList.remove("active"));
      button.classList.add("active");
      const mode = button.getAttribute("data-mode");
      if (mode === "side-by-side") {
        previewSideBySide.style.display = "flex";
        previewBoth.style.display = "none";
      } else if (mode === "both") {
        previewBoth.style.display = "flex";
        previewSideBySide.style.display = "none";
      }
    });
  });

  // --- Render Indicator on Active Tab ---
  function setRenderIndicator(show) {
    const activeTab = document.querySelector("#preview-tabs .tab-button.active");
    if (!activeTab) return;
    let indicator = activeTab.querySelector(".render-indicator");
    if (show) {
      if (!indicator) {
        indicator = document.createElement("span");
        indicator.className = "render-indicator";
        activeTab.appendChild(indicator);
      }
    } else {
      if (indicator) indicator.remove();
    }
  }

  // --- Helper: Auto-Render Timer Update ---
  function updateAutoRenderTimer() {
    if (!autoRenderScheduledTime) {
      autoRenderTimerDisplay.textContent = "";
      return;
    }
    const secondsLeft = Math.ceil((autoRenderScheduledTime - Date.now()) / 1000);
    autoRenderTimerDisplay.textContent = `Auto render in ${secondsLeft}s`;
  }
  setInterval(updateAutoRenderTimer, 1000);

  // --- Initialize Splits ---
  Split(['#editors-container', '#preview-container'], {
    sizes: [30, 70],
    minSize: 200,
    gutterSize: 6,
  });
  Split(
    Array.from(document.querySelectorAll('#editors-container .editor-section')),
    {
      direction: 'vertical',
      sizes: [25, 25, 25, 25],
      gutterSize: 6,
      minSize: 50,
    }
  );

  // --- Render Function: Only If Content Changed ---
  async function renderGFXS(force = false) {
    const currentContent = editorVars.getValue() + editorFilters.getValue() + editorComposition.getValue() + editorLayers.getValue();
    if (!force && currentContent === lastRenderedContent) return;
    if (isRenderInProgress && !force) return;
    isRenderInProgress = true;
    setRenderIndicator(true);
    const gfxs = "[VARS]\n" + editorVars.getValue() +
      "\n\n[FILTERS]\n" + editorFilters.getValue() +
      "\n\n[COMPOSITION]\n" + editorComposition.getValue() +
      "\n\n[LAYERS]\n" + editorLayers.getValue() + "\n";
  
    try {
      const response = await fetch("/render", {
        method: "POST",
        headers: { "Content-Type": "application/x-www-form-urlencoded" },
        body: `gfxs=${encodeURIComponent(gfxs)}`,
      });
      if (response.ok) {
        const json = await response.json();
        if (json.update) {
          // Store base64 data for later toggling in both mode.
          base64Original = json.original;
          base64Processed = json.processed;
          // Set the image sources using data URLs.
          originalImageSB.src = "data:image/png;base64," + json.original;
          processedImageSB.src = "data:image/png;base64," + json.processed;
          previewImage.src = "data:image/png;base64," + json.processed;
          lastRenderedContent = currentContent;
        } else {
          alert("No update from render.");
        }
      } else {
        alert("Failed to render: " + await response.text());
      }
    } catch (error) {
      alert("An error occurred: " + error.message);
    } finally {
      isRenderInProgress = false;
      setRenderIndicator(false);
      scheduleAutoRender();
    }
  }
  

  // --- Auto-Render Scheduling ---
  function scheduleAutoRender() {
    clearTimeout(autoRenderTimeoutId);
    if (!isRenderInProgress) {
      autoRenderScheduledTime = Date.now() + AUTO_RENDER_DELAY;
      autoRenderTimeoutId = setTimeout(renderGFXS, AUTO_RENDER_DELAY);
    }
  }

  // Reset auto-render timer on CodeMirror activity.
  [editorComposition, editorFilters, editorLayers, editorVars].forEach(editor => {
    editor.on("change", scheduleAutoRender);
    editor.on("cursorActivity", scheduleAutoRender);
  });

  // --- Force Render Button ---
  forceRenderBtn.addEventListener("click", () => {
    clearTimeout(autoRenderTimeoutId);
    autoRenderScheduledTime = null;
    updateAutoRenderTimer();
    renderGFXS(true);
  });

  // --- Filter Management Events ---
  filterDropdown.addEventListener("change", () => {
    const name = filterDropdown.value;
    if (name) {
      loadFilterByName(name);
      renderGFXS(true);
    }
  });
  saveFilterBtn.addEventListener("click", () => {
    saveCurrentFilter();
  });

  // --- Load Available Filters ---
  async function loadFilterList() {
    try {
      const response = await fetch("/filters");
      if (response.ok) {
        const names = await response.json();
        filterDropdown.innerHTML = '<option value="">-- Select a filter --</option>';
        names.forEach(name => {
          const option = document.createElement("option");
          option.value = name;
          option.textContent = name;
          filterDropdown.appendChild(option);
        });
      }
    } catch (err) {
      console.error("Failed to load filter list", err);
    }
  }

  // --- Load a Filter by Name ---
  async function loadFilterByName(name) {
    try {
      const response = await fetch("/filter?name=" + encodeURIComponent(name));
      if (response.ok) {
        const filterText = await response.text();
        const sections = {};
        const pattern = /\s*\[(VARS|FILTERS|COMPOSITION|LAYERS)\]\s*\n([\s\S]*?)(?=\n\s*\[(?:VARS|FILTERS|COMPOSITION|LAYERS)\]\s*\n|$)/g;
        let match;
        while ((match = pattern.exec(filterText)) !== null) {
          sections[match[1]] = match[2].trim();
        }
        if (sections.VARS) editorVars.setValue(sections.VARS);
        if (sections.FILTERS) editorFilters.setValue(sections.FILTERS);
        if (sections.COMPOSITION) editorComposition.setValue(sections.COMPOSITION);
        if (sections.LAYERS) editorLayers.setValue(sections.LAYERS);
      } else {
        alert("Failed to load filter");
      }
    } catch (err) {
      alert("Error loading filter: " + err.message);
    }
  }

  // --- Save the Current Filter ---
  async function saveCurrentFilter() {
    const filterText = "[VARS]\n" + editorVars.getValue() +
      "\n\n[FILTERS]\n" + editorFilters.getValue() +
      "\n\n[COMPOSITION]\n" + editorComposition.getValue() +
      "\n\n[LAYERS]\n" + editorLayers.getValue() + "\n";
    const name = prompt("Enter a name to save the filter:");
    if (!name) return;
    try {
      const formData = new FormData();
      formData.append("name", name);
      formData.append("filter", filterText);
      const response = await fetch("/saveFilter", {
        method: "POST",
        body: formData,
      });
      if (response.ok) {
        loadFilterList();
      } else {
        alert("Failed to save filter: " + await response.text());
      }
    } catch (err) {
      alert("Error saving filter: " + err.message);
    }
  }

  // --- Drag-and-Drop Handling for Image Upload / Batch Processing ---
  document.addEventListener("dragover", (e) => {
    e.preventDefault();
  });
  document.addEventListener("drop", async (e) => {
    e.preventDefault();
    if (e.dataTransfer.files.length === 0) return;
    if (e.dataTransfer.files.length === 1) {
      const file = e.dataTransfer.files[0];
      if (!file.type.startsWith("image/")) return;
      // For single image upload, we use a simple overlay (full-screen is still used for batch).
      showOverlay(overlayFull);
      const formData = new FormData();
      formData.append("image", file);
      try {
        const response = await fetch("/upload", {
          method: "POST",
          body: formData,
        });
        if (response.ok) {
          renderGFXS(true);
        } else {
          alert("Failed to upload image: " + await response.text());
        }
      } catch (err) {
        alert("An error occurred during image upload: " + err.message);
      } finally {
        hideOverlay(overlayFull);
      }
    } else {
      // Batch processing.
      showOverlay(overlayFull);
      const gfxs = "[VARS]\n" + editorVars.getValue() +
        "\n\n[FILTERS]\n" + editorFilters.getValue() +
        "\n\n[COMPOSITION]\n" + editorComposition.getValue() +
        "\n\n[LAYERS]\n" + editorLayers.getValue() + "\n";
      const formData = new FormData();
      formData.append("gfxs", gfxs);
      for (let i = 0; i < e.dataTransfer.files.length; i++) {
        formData.append("images", e.dataTransfer.files[i]);
      }
      try {
        const response = await fetch("/renderBatch", {
          method: "POST",
          body: formData,
        });
        if (response.ok) {
          const blob = await response.blob();
          const url = URL.createObjectURL(blob);
          const a = document.createElement("a");
          a.href = url;
          a.download = "batch.zip";
          document.body.appendChild(a);
          a.click();
          a.remove();
          URL.revokeObjectURL(url);
        } else {
          alert("Batch render failed: " + await response.text());
        }
      } catch (err) {
        alert("An error occurred during batch render: " + err.message);
      } finally {
        hideOverlay(overlayFull);
        renderGFXS(true);
      }
    }
  });

  // --- Initial Load ---
  loadFilterList();
  loadFilterByName("default");
  renderGFXS();
  
  // --- New Feature: Pixel Info on Hover ---
  const imageElements = [previewImage, originalImageSB, processedImageSB];
  // Map to store offscreen canvases for each image element.
  const imageCanvasMap = new WeakMap();

  // Function to update/create offscreen canvas for an image.
  function updateCanvasForImage(img) {
    if (!img.naturalWidth || !img.naturalHeight) return;
    const canvas = document.createElement("canvas");
    canvas.width = img.naturalWidth;
    canvas.height = img.naturalHeight;
    const ctx = canvas.getContext("2d");
    ctx.drawImage(img, 0, 0);
    imageCanvasMap.set(img, canvas);
  }

  // When an image loads, update its offscreen canvas.
  imageElements.forEach(img => {
    img.addEventListener("load", () => {
      updateCanvasForImage(img);
    });
    // In case image is already loaded:
    if (img.complete) updateCanvasForImage(img);
  });

  // Create (or get) a pixel info overlay for an image.
  function getPixelInfoOverlay(img) {
    let overlay = img.parentElement.querySelector(".pixel-info");
    if (!overlay) {
      overlay = document.createElement("div");
      overlay.className = "pixel-info";
      // Initially hidden.
      overlay.style.display = "none";
      img.parentElement.appendChild(overlay);
    }
    return overlay;
  }

  // Convert RGB to HSL
  function rgbToHsl(r, g, b) {
    r /= 255;
    g /= 255;
    b /= 255;
    const max = Math.max(r, g, b), min = Math.min(r, g, b);
    let h, s, l = (max + min) / 2;
    if(max === min) {
      h = s = 0; // achromatic
    } else {
      const d = max - min;
      s = l > 0.5 ? d / (2 - max - min) : d / (max + min);
      switch(max) {
        case r: h = (g - b) / d + (g < b ? 6 : 0); break;
        case g: h = (b - r) / d + 2; break;
        case b: h = (r - g) / d + 4; break;
      }
      h *= 60;
    }
    return { h: Math.round(h), s: Math.round(s * 100), l: Math.round(l * 100) };
  }

  // Mouse move handler to show pixel info.
  function handleMouseMove(e) {
    const img = e.target;
    const rect = img.getBoundingClientRect();
    // Calculate offset relative to the image.
    const offsetX = e.clientX - rect.left;
    const offsetY = e.clientY - rect.top;
    // Map to actual image coordinates.
    const actualX = Math.floor(offsetX * img.naturalWidth / img.clientWidth);
    const actualY = Math.floor(offsetY * img.naturalHeight / img.clientHeight);
    
    const canvas = imageCanvasMap.get(img);
    if (!canvas) return;
    const ctx = canvas.getContext("2d");
    // Ensure coordinates are within bounds.
    if (actualX < 0 || actualY < 0 || actualX >= canvas.width || actualY >= canvas.height) return;
    const pixelData = ctx.getImageData(actualX, actualY, 1, 1).data;
    const [r, g, b, a] = pixelData;
    const hsla = rgbToHsl(r, g, b);
    const alpha = (a / 255).toFixed(2);
    
    // Build the info text.
    const infoText = `x: ${actualX}, y: ${actualY} | H: ${hsla.h}°, S: ${hsla.s}%, L: ${hsla.l}%, A: ${alpha}`;
    
    // Get or create overlay.
    const overlay = getPixelInfoOverlay(img);
    // Create or update color circle.
    let colorCircle = overlay.querySelector(".color-circle");
    if (!colorCircle) {
      colorCircle = document.createElement("div");
      colorCircle.className = "color-circle";
      overlay.insertBefore(colorCircle, overlay.firstChild);
    }
    // Set the circle background to the pixel color.
    colorCircle.style.backgroundColor = `rgba(${r}, ${g}, ${b}, ${alpha})`;
    
    overlay.textContent = ""; // Clear previous text
    overlay.appendChild(colorCircle);
    const textNode = document.createTextNode(infoText);
    overlay.appendChild(textNode);
    
    overlay.style.display = "flex";
  }

  // Mouse leave handler to hide pixel info.
  function handleMouseLeave(e) {
    const overlay = getPixelInfoOverlay(e.target);
    overlay.style.display = "none";
  }

  // Attach the event handlers to all images.
  imageElements.forEach(img => {
    img.addEventListener("mousemove", handleMouseMove);
    img.addEventListener("mouseleave", handleMouseLeave);
  });

});
