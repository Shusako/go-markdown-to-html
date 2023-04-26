document.addEventListener('DOMContentLoaded', function () {
    let files = [];

    fetch("files.json")
        .then((response) => response.json())
        .then((data) => {
            files = data;
            renderFileList();
        })
        .catch((error) => {
            console.error("Error fetching files.json:", error);
        });

    function renderFileList() {
        const fileList = document.getElementById("fileList");

        files.forEach((file) => {
            const listItem = document.createElement("li");
            const link = document.createElement("a");

            link.href = file;
            link.textContent = file;

            listItem.appendChild(link);
            fileList.appendChild(listItem);

            listItem.addEventListener("click", function (event) {
                event.preventDefault(); // Prevent the default link behavior
                loadFile(file);
            });
        });
    }

    function loadFile(file) {
        fetch(file)
            .then((response) => response.text())
            .then((html) => {
                const content = document.getElementById("content");
                content.innerHTML = html;
            })
            .catch((error) => {
                console.error("Error loading file:", error);
            });
    }

});
