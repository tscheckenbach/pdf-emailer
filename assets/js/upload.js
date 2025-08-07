document.addEventListener('DOMContentLoaded', function() {
    const dropZone = document.getElementById('dropZone');
    const fileInput = document.getElementById('fileInput');
    const fileList = document.getElementById('fileList');
    const fileItems = document.getElementById('fileItems');
    const uploadForm = document.getElementById('uploadForm');
    const uploadBtn = document.getElementById('uploadBtn');
    const uploadProgress = document.getElementById('uploadProgress');
    const progressBar = document.getElementById('progressBar');
    const progressText = document.getElementById('progressText');
    
    let selectedFiles = [];

    // Drag and drop functionality
    dropZone.addEventListener('dragover', function(e) {
        e.preventDefault();
        dropZone.classList.add('border-blue-400', 'bg-blue-50');
    });

    dropZone.addEventListener('dragleave', function(e) {
        e.preventDefault();
        dropZone.classList.remove('border-blue-400', 'bg-blue-50');
    });

    dropZone.addEventListener('drop', function(e) {
        e.preventDefault();
        dropZone.classList.remove('border-blue-400', 'bg-blue-50');
        
        const files = Array.from(e.dataTransfer.files).filter(file => file.type === 'application/pdf');
        addFiles(files);
    });

    dropZone.addEventListener('click', function() {
        fileInput.click();
    });

    fileInput.addEventListener('change', function(e) {
        const files = Array.from(e.target.files);
        addFiles(files);
    });

    function addFiles(files) {
        files.forEach(file => {
            if (file.type === 'application/pdf' && !selectedFiles.find(f => f.name === file.name)) {
                selectedFiles.push(file);
            }
        });
        updateFileList();
    }

    function updateFileList() {
        if (selectedFiles.length > 0) {
            fileList.classList.remove('hidden');
            fileItems.innerHTML = '';
            
            selectedFiles.forEach((file, index) => {
                const fileItem = document.createElement('div');
                fileItem.className = 'flex items-center justify-between p-3 bg-gray-50 rounded-lg';
                fileItem.innerHTML = `
                    <div class="flex items-center space-x-3">
                        <svg class="w-5 h-5 text-red-500" fill="currentColor" viewBox="0 0 20 20">
                            <path fill-rule="evenodd" d="M4 4a2 2 0 00-2 2v8a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2H4zm0 2h12v8H4V6z" clip-rule="evenodd"></path>
                        </svg>
                        <span class="text-sm font-medium text-gray-900">${file.name}</span>
                        <span class="text-xs text-gray-500">(${(file.size / 1024 / 1024).toFixed(2)} MB)</span>
                    </div>
                    <button type="button" onclick="removeFile(${index})" class="text-red-500 hover:text-red-700">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                        </svg>
                    </button>
                `;
                fileItems.appendChild(fileItem);
            });
        } else {
            fileList.classList.add('hidden');
        }
    }

    window.removeFile = function(index) {
        selectedFiles.splice(index, 1);
        updateFileList();
    };

    // Form submission
    uploadForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const template = document.getElementById('template').value;
        if (!template) {
            alert('Es muss ein Template ausgewählt werden');
            return;
        }
        
        if (selectedFiles.length === 0) {
            alert('Es müssen PDF-Dateien hinzugefügt werden');
            return;
        }

        uploadBtn.disabled = true;
        uploadProgress.classList.remove('hidden');
        
        let completedUploads = 0;
        const totalFiles = selectedFiles.length;
        
        for (let i = 0; i < selectedFiles.length; i++) {
            const file = selectedFiles[i];
            const formData = new FormData();
            formData.append('file', file);
            formData.append('template', template);
            
            try {
                const response = await fetch('/upload', {
                    method: 'POST',
                    body: formData
                });
                
                if (response.ok) {
                    console.log(`File ${file.name} uploaded successfully`);
                } else {
                    console.error(`Error uploading ${file.name}:`, await response.text());
                }
            } catch (error) {
                console.error(`Error uploading ${file.name}:`, error);
            }
            
            completedUploads++;
            const progress = Math.round((completedUploads / totalFiles) * 100);
            progressBar.style.width = progress + '%';
            progressText.textContent = progress + '%';
        }
        
        // Reset form after upload
        setTimeout(() => {
            selectedFiles = [];
            updateFileList();
            uploadForm.reset();
            uploadBtn.disabled = false;
            uploadProgress.classList.add('hidden');
            progressBar.style.width = '0%';
            progressText.textContent = '0%';
            alert('Upload completed!');
        }, 1000);
    });
});
