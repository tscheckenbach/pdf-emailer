document.addEventListener('DOMContentLoaded', function() {
    const dropZone = document.getElementById('dropZone');
    const fileInput = document.getElementById('fileInput');
    const fileListContainer = document.getElementById('fileListContainer');
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
            fileListContainer.innerHTML = `
                <h3 class="text-sm font-medium text-gray-700 mb-3">Selected Files:</h3>
                <div class="space-y-2">
                    ${selectedFiles.map((file, index) => `
                        <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg" data-filename="${file.name}">
                            <div class="flex items-center space-x-3">
                                <div class="status-icon">
                                    <svg class="w-5 h-5 text-gray-400" fill="currentColor" viewBox="0 0 20 20">
                                        <path fill-rule="evenodd" d="M4 4a2 2 0 00-2 2v8a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2H4zm0 2h12v8H4V6z" clip-rule="evenodd"></path>
                                    </svg>
                                </div>
                                <span class="text-sm font-medium text-gray-900">${file.name}</span>
                                <span class="text-xs text-gray-500">(${(file.size / 1024 / 1024).toFixed(2)} MB)</span>
                            </div>
                            <button type="button" onclick="removeFile(${index})" class="text-red-500 hover:text-red-700">
                                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                                </svg>
                            </button>
                        </div>
                    `).join('')}
                </div>
            `;
        } else {
            fileListContainer.innerHTML = '';
        }
    }

    window.removeFile = function(index) {
        selectedFiles.splice(index, 1);
        updateFileList();
    };

    window.uploadFiles = async function() {
        const template = document.getElementById('template').value;
        if (!template) {
        alert('Es muss eine Vorlage ausgewählt werden');
            return;
        }
        
        if (selectedFiles.length === 0) {
        alert('Es müssen PDF-Dateien hinzugefügt werden');
            return;
        }

        uploadProgress.classList.remove('hidden');

        // Small helpers (vanilla JS only)
        const spinnerHtml = `
            <svg class="w-5 h-5 text-blue-500 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
            </svg>
        `;
        const setStatusIcon = (fileName, html) => {
            const fileItem = document.querySelector(`[data-filename="${fileName}"]`);
            if (!fileItem) return;
            const statusIcon = fileItem.querySelector('.status-icon');
            if (statusIcon) statusIcon.innerHTML = html;
        };
        const networkErrorHtml = (message) => `
            <div class="relative group">
                <svg class="w-5 h-5 text-red-500 cursor-help" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"></path>
                </svg>
                <div class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-3 py-2 bg-red-600 text-white text-sm rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none whitespace-nowrap z-10">
                    Netzwerkfehler: ${message}
                    <div class="absolute top-full left-1/2 transform -translate-x-1/2 border-4 border-transparent border-t-red-600"></div>
                </div>
            </div>
        `;
        const updateProgress = (completed, total) => {
            const progress = Math.round((completed / total) * 100);
            progressBar.style.width = progress + '%';
            progressText.textContent = progress + '%';
        };

        // Show loading state for all files
        selectedFiles.forEach((file) => setStatusIcon(file.name, spinnerHtml));
        
        const filesToUpload = [...selectedFiles];
        let completedUploads = 0;
        const totalFiles = filesToUpload.length;
        const MAX_CONCURRENCY = 5;

        const uploadSingle = async (file) => {
            const formData = new FormData();
            formData.append('file', file);
            formData.append('template', template);

            try {
                const response = await fetch('/upload', {
                    method: 'POST',
                    body: formData,
                    headers: {
                        'HX-Request': 'true'
                    }
                });

                const html = await response.text();
                setStatusIcon(file.name, html);
            } catch (error) {
                setStatusIcon(file.name, networkErrorHtml(error.message));
            } finally {
                completedUploads++;
                updateProgress(completedUploads, totalFiles);
            }
        };

        for (let i = 0; i < filesToUpload.length; i += MAX_CONCURRENCY) {
            const chunk = filesToUpload.slice(i, i + MAX_CONCURRENCY);
            await Promise.all(chunk.map(file => uploadSingle(file)));
        }
        
        setTimeout(() => {
            uploadProgress.classList.add('hidden');
            progressBar.style.width = '0%';
            progressText.textContent = '0%';
            
            // Replace the upload button with "Neuer Versand" button
            const uploadButtonContainer = document.getElementById('uploadButtonContainer');
            if (uploadButtonContainer) {
                uploadButtonContainer.innerHTML = `
                    <div class="flex justify-center">
                        <a href="/" class="inline-flex items-center px-6 py-3 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500">
                            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
                            </svg>
                            Neuer Versand
                        </a>
                    </div>
                `;
            }
        }, 2000);
    };
});
