/*
 * static/js/app.js - Main application logic
 */
document.addEventListener('DOMContentLoaded', () => {
    const toggleFormButton = document.getElementById('toggleFormButton');
    const rightPanel = document.getElementById('rightPanel');
    const markdownFileList = document.getElementById('markdownFileList');
    const newFileForm = document.getElementById('newFileForm');
    const filePreview = document.getElementById('filePreview');
    const previewFileName = document.getElementById('previewFileName');
    const previewContent = document.getElementById('previewContent');
    const checkGrammarButton = document.getElementById('checkGrammarButton');

    const customModal = document.getElementById('customModal');
    const modalTitle = document.getElementById('modalTitle');
    const modalMessage = document.getElementById('modalMessage');
    const modalCloseButton = document.getElementById('modalCloseButton');

    let isPanelActive = false; // To track the state of the right panel
    let currentPanelMode = 'form'; // 'form' or 'preview'
    let currentlyPreviewedFile = null; // To keep track of the file being viewed
    let selectedFileContents = '';

    // --- Custom Modal Functions ---
    /**
     * Displays a custom modal with a title and message.
     * @param {string} title - The title for the modal.
     * @param {string} message - The message content for the modal.
     */
    function showModal(title, message) {
        modalTitle.textContent = title;
        modalMessage.innerHTML = message;
        customModal.classList.remove('hidden');
    }

    /**
     * Hides the custom modal.
     */
    function hideModal() {
        customModal.classList.add('hidden');
    }

    modalCloseButton.addEventListener('click', hideModal);
    // Close modal if overlay is clicked (optional, but good for UX)
    customModal.addEventListener('click', (event) => {
        if (event.target === customModal) {
            hideModal();
        }
    });

    // --- Panel Visibility and Content Management ---

    /**
     * Toggles the visibility of the right panel and switches its content to the form.
     */
    function togglePanel() {
        isPanelActive = !isPanelActive;
        if (isPanelActive) {
            rightPanel.classList.add('active');
            showForm(); // Always show form when panel is activated by button
        } else {
            rightPanel.classList.remove('active');
            // Reset content when hidden
            newFileForm.classList.remove('hidden');
            filePreview.classList.add('hidden');
        }
    }

    /**
     * Displays the form within the right panel.
     */
    function showForm() {
        currentPanelMode = 'form';
        newFileForm.classList.remove('hidden');
        filePreview.classList.add('hidden');
    }

    /**
     * Displays the file preview within the right panel.
     * @param {string} fileName - The name of the file being previewed.
     * @param {string} content - The content of the file.
     */
    function showFilePreview(fileName, content) {
        currentPanelMode = 'preview';
        currentlyPreviewedFile = fileName; // Set the currently previewed file
        newFileForm.classList.add('hidden');
        filePreview.classList.remove('hidden');
        previewFileName.textContent = fileName;
        previewContent.innerHTML = content;
        selectedFileContents = content;
    }

    // --- Event Listeners ---

    toggleFormButton.addEventListener('click', togglePanel);

    // form submission
    newFileForm.addEventListener('submit', async (event) => {
        event.preventDefault(); // Prevent default form submission

        const fileNote = document.getElementById('fileNote');

        if (!fileNote.files.length) {
            showModal('Input Error', 'Please select a file before submitting.');
            return;
        }

        const file = fileNote.files[0];
        const formData = new FormData();
        formData.append('note', file);

        await fetch('/api/notes/save', {
                method: 'POST',
                body: formData,
                headers: {
                    'Accept': 'application/json',
                },
            })
            .then(response => {
                if (!response.ok) {
                    throw response
                }
                return response.json();
            })
            .then(data => {
                renderMarkdownFileList(); // Re-render the list to show the new file
                newFileForm.reset(); // Clear the form

                showModal('Success', data.message);
            })
            .catch(error => {
                console.log("Error occurred");
                try {
                    error.json().then(body => {
                        showModal(error.status, body.message)
                    });
                } catch (e) {
                    console.log("Error parsing promise");
                    console.error(error);
                }
                return
            });
    });

    /**
     * Handles click events on the delete button for a markdown file.
     * @param {Event} event - The click event.
     */
    async function handleDeleteFile(event) {
        event.stopPropagation(); // Prevent the parent <li> from triggering handleFileClick
        const button = event.currentTarget;
        const fileName = button.dataset.fileName;

        const confirmDelete = confirm(`Are you sure you want to delete "${fileName}"?`); // A simple confirm for demo
        if (!confirmDelete) {
            return;
        }

        // Get parent li
        const parentLi = button.parentElement

        const apiUrl = `/api/notes/${fileName}`;

        return await fetch(apiUrl, {
                headers: {
                    'Content-Type': 'application/json'
                },
                method: 'DELETE',
            })
            .then(response => response.json())
            .then(response => {
                parentLi.classList.add('removing');

                parentLi.style.opacity = 0; // Start fade out

                parentLi.addEventListener('transitionend', function() {
                    parentLi.remove();
                }, {
                    once: true
                }); // remove listener after execution
            })
            .catch(error => console.error('Error deleting file:', error));
    }

    /**
     * A grammar check API 
     * @param {string} text - The text to check for grammar.
     */
    async function checkGrammar(text) {
        showModal('Checking Grammar', 'Please wait while we check the grammar...');
        try {
            const payload = {
                text: text,
                language: 'en-US'
            };

            const result = await fetch('/api/notes/check-grammers', {
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    method: 'POST',
                    body: JSON.stringify(payload)
                })
                .then(response => response.json())
                .then(response => {
                    return response.data
                })
                .catch(error => console.error('Error fetching files:', error));

            let modifiedContent = []
            if (result.matches) {
                result.matches.forEach(match => {
                    const key = match.replacements[0] ? match.replacements[0].value : ''
                    const value = match.replacements[1] ? match.replacements[1].value : ''
                    if (key && value) {
                        modifiedContent.push({
                            original: value,
                            suggestion: key
                        });
                    }
                });
            }

            // Display the modified content
            if (modifiedContent.length > 0) {
                let tblContent = '<table class="table-auto">  <thead>    <tr>      <th>Original</th>      <th>Suggestion</th>    </tr>  </thead>  <tbody>'
                modifiedContent.forEach(item => {
                    tblContent += `    <tr>      <td>${item.original}</td>      <td>${item.suggestion}</td>    </tr>`
                })
                tblContent += '  </tbody></table>'
                showModal('Grammar Check Results', tblContent);
            } else {
                showModal('Grammar Check Failed', 'Could not get grammar check results. Please try again.');
            }
        } catch (error) {
            console.error('Error during grammar check API call:', error);
            showModal('Grammar Check Error', 'An error occurred while checking grammar. Please try again.');
        }
    }

    checkGrammarButton.addEventListener('click', () => {
        checkGrammar(selectedFileContents);
    });


    // --- Data Fetching and Rendering ---

    /**
     * Fetching a list of markdown files from an API.
     * @returns {Promise<Array<string>>} A promise that resolves with an array of file names.
     */
    async function fetchMarkdownFiles() {
        return await fetch('/api/notes', {
                headers: {
                    'Content-Type': 'application/json'
                },
                method: 'GET',
            })
            .then(response => response.json())
            .then(response => {
                return response.data
            })
            .catch(error => console.error('Error fetching files:', error));
    }

    /**
     * Fetching the content of a specific markdown file from an API.
     * @param {string} fileName - The name of the file to fetch.
     * @returns {Promise<string>} A promise that resolves with the file content.
     */
    async function fetchFileContent(fileName) {
        return await fetch(`/api/notes/${fileName}`, {
                headers: {
                    'Content-Type': 'application/json'
                },
                method: 'GET',
            })
            .then(response => response.json())
            .then(response => {
                return response.data
            })
            .catch(error => console.error('Error fetching file:', error));
    }

    /**
     * Renders the list of markdown files in the left column.
     */
    async function renderMarkdownFileList() {
        try {
            const files = await fetchMarkdownFiles();

            markdownFileList.innerHTML = ''; // Clear previous list
            files.forEach(file => {
                const fileTitle = file.title
                const fileName = file.file_name
                const fileUrl = file.url

                const li = document.createElement('li');
                // add id to li to remove if deleted
                li.id = fileName
                // Add 'active' class if this is the currently previewed file
                const isActive = (currentlyPreviewedFile === fileName) ? 'bg-blue-100 border-blue-300' : '';
                li.className = `cursor-pointer p-3 bg-white hover:bg-blue-50 transition-colors duration-200 rounded-md shadow-sm border border-gray-200 flex items-center justify-between transition-opacity duration-300 ease-in-out ${isActive}`;
                li.innerHTML = `
                    <span class="text-blue-700 font-medium flex-grow">${fileName}</span>
                    <button class="delete-file-button ml-2 p-1 text-red-500 hover:text-red-700 rounded-full focus:outline-none focus:ring-2 focus:ring-red-500" data-file-name="${fileName}" aria-label="Delete ${fileName}">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                            <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 011-1h4a1 1 0 110 2H8a1 1 0 01-1-1zm2 3a1 1 0 011-1h0a1 1 0 110 2h0a1 1 0 01-1-1zm2 3a1 1 0 011-1h0a1 1 0 110 2h0a1 1 0 01-1-1z" clip-rule="evenodd" />
                        </svg>
                    </button>
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-400 ml-2" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                    </svg>
                `;
                li.dataset.fileName = fileName; // Store file name for easy access
                li.addEventListener('click', handleFileClick);

                // Attach event listener to the delete button specifically
                const deleteButton = li.querySelector('.delete-file-button');
                if (deleteButton) {
                    deleteButton.addEventListener('click', handleDeleteFile);
                }
                markdownFileList.appendChild(li);
            });
        } catch (error) {
            console.error('Failed to load markdown files:', error);
            markdownFileList.innerHTML = '<li class="text-red-500 p-3">Error loading files. Please try again.</li>';
        }
    }

    /**
     * Handles click events on markdown file list items.
     * @param {Event} event - The click event.
     */
    async function handleFileClick(event) {
        const clickedElement = event.currentTarget;
        const fileName = clickedElement.dataset.fileName;

        // Remove active class from previously selected item
        const activeItem = markdownFileList.querySelector('.bg-blue-100');
        if (activeItem) {
            activeItem.classList.remove('bg-blue-100', 'border-blue-300');
        }
        // Add active class to the clicked item
        clickedElement.classList.add('bg-blue-100', 'border-blue-300');


        try {
            // Show loading state
            previewFileName.textContent = fileName;
            previewContent.textContent = 'Loading file content...';
            if (!isPanelActive) {
                togglePanel(); // Show panel if hidden
            }
            showFilePreview(fileName, 'Loading file content...'); // Show preview container

            const content = await fetchFileContent(fileName);
            showFilePreview(fileName, content); // Update with actual content
        } catch (error) {
            console.error(`Failed to load content for ${fileName}:`, error);
            showFilePreview(fileName, `Error: Could not load content for "${fileName}".`);
        }
    }

    // Initial render when the page loads
    renderMarkdownFileList();
});