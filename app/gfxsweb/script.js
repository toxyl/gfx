document.addEventListener("DOMContentLoaded", () => {
    const editor = CodeMirror.fromTextArea(document.getElementById("gfxs-editor"), {
        lineNumbers: true,
        mode: "gfxscript",
    });

    const renderBtn = document.getElementById("render-btn");
    const renderedImage = document.getElementById("rendered-image");

    renderBtn.addEventListener("click", async () => {
        const gfxs = editor.getValue();

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
            } else {
                alert("Failed to render: " + await response.text());
            }
        } catch (error) {
            alert("An error occurred: " + error.message);
        }
    });
});
