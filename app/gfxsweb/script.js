document.addEventListener("DOMContentLoaded", () => {
    const editor = CodeMirror.fromTextArea(document.getElementById("gfxs-editor"), {
        lineNumbers: true,
        mode: "gfxscript",
    });

    const renderedImage = document.getElementById("rendered-image");
    const imageContainer = document.getElementById("image-container"); // Updated target
    const outputDiv = document.getElementById("output");
    let timeoutId = null;

    const renderGFXS = async () => {
        const gfxs = editor.getValue();

        // Show a loading message
        outputDiv.classList.add("loading");
        imageContainer.innerHTML = "<p>&nbsp;Rendering...&nbsp;</p>"; // Update within the container

        try {
            const response = await fetch("/render", {
                method: "POST",
                headers: {
                    "Content-Type": "application/x-www-form-urlencoded",
                },
                body: `gfxs=${encodeURIComponent(gfxs)}`,
            });

            if (response.ok) {
                const blob = await response.blob();
                renderedImage.src = URL.createObjectURL(blob);
                renderedImage.onload = () => URL.revokeObjectURL(renderedImage.src); // Clean up memory

                imageContainer.innerHTML = ""; // Clear the loading message
                imageContainer.appendChild(renderedImage); // Add the rendered image back
            } else {
                alert("Failed to render: " + await response.text());
                imageContainer.innerHTML = ""; // Clear the loading message on error
            }
        } catch (error) {
            alert("An error occurred: " + error.message);
            imageContainer.innerHTML = ""; // Clear the loading message on error
        } finally {
            outputDiv.classList.remove("loading");
        }
    };

    editor.on("change", () => {
        clearTimeout(timeoutId); // Reset the timer on every keypress
        timeoutId = setTimeout(renderGFXS, 5000); // Trigger render after 5 seconds of inactivity
    });

    renderGFXS();
});
