document.addEventListener("DOMContentLoaded", function () {
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
        const parser = new DOMParser();
        const fetchedDoc = parser.parseFromString(html, "text/html");
        const fetchedContent = fetchedDoc.querySelector("#content");

        if (fetchedContent) {
          const content = document.getElementById("content");
          content.innerHTML = fetchedContent.innerHTML;
        } else {
          console.error("No content element found in the fetched file:", file);
        }
      })
      .catch((error) => {
        console.error("Error loading file:", error);
      });
  }

  document
    .getElementById("toggle-file-list")
    .addEventListener("click", function () {
      var fileList = document.querySelector(".file-list");
      var overlay = document.getElementById("overlay");

      fileList.classList.toggle("show");
      overlay.classList.toggle("show");
    });

  function updateFileListDisplay() {
    var fileList = document.querySelector(".file-list");
    var toggleFileListButton = document.getElementById("toggle-file-list");

    if (window.innerWidth <= 1288) {
      toggleFileListButton.style.display = "block";
    } else {
      fileList.classList.remove("show");
      toggleFileListButton.style.display = "none";
    }
  }

  window.addEventListener("resize", updateFileListDisplay);
  updateFileListDisplay();
});
