/**
 * Image Metadata Viewer - Main JavaScript
 */

// Tab switching functionality
function switchTab(tabName) {
    // Hide all tabs
    document.querySelectorAll('.tab-content').forEach(tab => {
        tab.classList.remove('active');
    });
    document.querySelectorAll('.tab').forEach(tab => {
        tab.classList.remove('active');
    });

    // Show selected tab
    document.getElementById(tabName + '-tab').classList.add('active');
    event.target.classList.add('active');
}

// File size formatter
function formatBytes(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}

// Display selected files
function displayFiles(files) {
    const fileList = document.getElementById('file-list');
    fileList.innerHTML = '';

    if (files.length === 0) return;

    for (let file of files) {
        const fileItem = document.createElement('div');
        fileItem.className = 'file-item';
        fileItem.innerHTML = `
            <span class="file-name">${file.name}</span>
            <span class="file-size">${formatBytes(file.size)}</span>
        `;
        fileList.appendChild(fileItem);
    }
}

// Initialize upload functionality when DOM is ready
document.addEventListener('DOMContentLoaded', function () {
    const uploadArea = document.getElementById('upload-area');
    const fileInput = document.getElementById('file-input');
    const fileList = document.getElementById('file-list');

    // Only initialize if upload elements exist (home page)
    if (uploadArea && fileInput) {
        // Click to upload
        uploadArea.addEventListener('click', () => fileInput.click());

        // Drag over effect
        uploadArea.addEventListener('dragover', (e) => {
            e.preventDefault();
            uploadArea.classList.add('drag-over');
        });

        // Drag leave effect
        uploadArea.addEventListener('dragleave', () => {
            uploadArea.classList.remove('drag-over');
        });

        // Drop handler
        uploadArea.addEventListener('drop', (e) => {
            e.preventDefault();
            uploadArea.classList.remove('drag-over');
            fileInput.files = e.dataTransfer.files;
            displayFiles(e.dataTransfer.files);
        });

        // File input change handler
        fileInput.addEventListener('change', (e) => {
            displayFiles(e.target.files);
        });
    }

    // Try API functionality (docs page)
    const responseBox = document.getElementById('try-response');
    const getForm = document.getElementById('try-get-form');
    const postForm = document.getElementById('try-post-form');

    if (responseBox && (getForm || postForm)) {
        const setResponse = (message, isError) => {
            responseBox.textContent = message;
            responseBox.classList.toggle('response-error', Boolean(isError));
        };

        const fetchAndRender = async (request) => {
            setResponse('Loading...', false);
            try {
                const res = await request;
                const contentType = res.headers.get('Content-Type') || '';
                const raw = await res.text();
                let formatted = raw;
                if (contentType.includes('application/json')) {
                    try {
                        formatted = JSON.stringify(JSON.parse(raw), null, 2);
                    } catch (err) {
                        formatted = raw;
                    }
                }
                setResponse(`Status: ${res.status} ${res.statusText}\n\n${formatted}`, !res.ok);
            } catch (err) {
                setResponse(`Request failed: ${err.message}`, true);
            }
        };

        if (getForm) {
            getForm.addEventListener('submit', (e) => {
                e.preventDefault();
                const input = document.getElementById('try-get-url');
                const url = input ? input.value.trim() : '';
                if (!url) {
                    setResponse('Please enter a URL.', true);
                    return;
                }
                if (url.includes('?') || url.includes('#')) {
                    setResponse('Use POST for URLs with query or hash fragments.', true);
                    return;
                }
                const safeUrl = encodeURI(url);
                fetchAndRender(fetch(`/api/${safeUrl}`));
            });
        }

        if (postForm) {
            postForm.addEventListener('submit', (e) => {
                e.preventDefault();
                const input = document.getElementById('try-post-urls');
                const raw = input ? input.value : '';
                const urls = raw
                    .split('\n')
                    .map((value) => value.trim())
                    .filter(Boolean);
                if (urls.length === 0) {
                    setResponse('Please enter at least one URL.', true);
                    return;
                }
                fetchAndRender(fetch('/api', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ urls }),
                }));
            });
        }
    }
});
