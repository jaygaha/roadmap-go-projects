document.addEventListener('DOMContentLoaded', () => {
    const tabs = document.querySelectorAll('[role="tab"]');
    const tabPanels = document.querySelectorAll('[role="tabpanel"]');
    
    // Initialize tab panels with proper transition styles
    // These styles are also defined in style.css but we ensure they're applied here as well
    tabPanels.forEach(panel => {
      // Apply transition properties
      panel.style.transition = 'opacity 0.3s ease-in-out, transform 0.3s ease-in-out';
      
      // Set initial state based on visibility
      if (panel.classList.contains('hidden')) {
        panel.style.opacity = '0';
        panel.style.transform = 'translateY(10px)';
        panel.style.pointerEvents = 'none';
      } else {
        panel.style.opacity = '1';
        panel.style.transform = 'translateY(0)';
        panel.style.pointerEvents = 'auto';
      }
    });

    // Global function for onclick handlers in HTML
    window.showTab = function(tabNumber) {
      // Convert from 1-based (HTML) to 0-based (JS array)
      const index = tabNumber - 1;
      activateTab(index);
      return false; // Prevent default action and stop propagation
    };

    function activateTab(index) {
      // First, start the hiding animation for all panels
      tabPanels.forEach(panel => {
        if (!panel.classList.contains('hidden')) {
          // Only animate currently visible panels
          panel.style.opacity = '0';
          panel.style.transform = 'translateY(10px)';
          panel.style.pointerEvents = 'none';
        }
      });
      
      // Reset all tab button styles
      tabs.forEach(tab => {
        tab.classList.remove('text-blue-600', 'border-b-2', 'border-blue-600', 'font-semibold');
        tab.classList.add('text-gray-500', 'border-transparent', 'font-medium');
        tab.setAttribute('aria-selected', 'false');
      });
  
      // Activate selected tab button
      tabs[index].classList.add('text-blue-600', 'border-b-2', 'border-blue-600', 'font-semibold');
      tabs[index].classList.remove('text-gray-500', 'border-transparent', 'font-medium');
      tabs[index].setAttribute('aria-selected', 'true');
      
      // Use setTimeout to create a smooth transition between panels
      setTimeout(() => {
        // Hide all panels
        tabPanels.forEach(panel => {
          panel.classList.add('hidden');
        });
        
        // Show selected panel
        const selectedPanel = tabPanels[index];
        selectedPanel.classList.remove('hidden');
        
        // Trigger reflow to ensure transition works
        selectedPanel.offsetHeight;
        
        // Apply visible state with transition
        selectedPanel.style.opacity = '1';
        selectedPanel.style.transform = 'translateY(0)';
        selectedPanel.style.pointerEvents = 'auto';
      }, 300); // Match this timing with the CSS transition duration
    }
  
    tabs.forEach((tab, i) => {
      tab.addEventListener('click', e => {
        e.preventDefault();
        activateTab(i);
      });
    });
  
    // Initialize first tab as active
    activateTab(0);
  
    // Keyboard navigation for accessibility
    document.addEventListener('keydown', e => {
      const activeIndex = Array.from(tabs).findIndex(t => t.getAttribute('aria-selected') === 'true');
      if (e.key === 'ArrowRight') {
        const next = (activeIndex + 1) % tabs.length;
        tabs[next].focus();
        activateTab(next);
        e.preventDefault();
      } else if (e.key === 'ArrowLeft') {
        const prev = (activeIndex - 1 + tabs.length) % tabs.length;
        tabs[prev].focus();
        activateTab(prev);
        e.preventDefault();
      }
    });
  
    function submitForm(form) {
      const formData = new FormData(form);
      const data = {};
      formData.forEach((value, key) => {
        data[key] = value;
      });
  
      // Determine API endpoint by form id or parent container id
      let apiEndpoint = '';
      let successElement, errorElement, formElement, resultElement, errorMessageElement;
      
      if (form.closest('#tab1')) {
        apiEndpoint = '/api/length';
        successElement = document.getElementById('length-success');
        errorElement = document.getElementById('length-error');
        formElement = document.getElementById('length-form');
        resultElement = document.getElementById('length-result');
        errorMessageElement = document.getElementById('length-error-message');
      } else if (form.closest('#tab2')) {
        apiEndpoint = '/api/weight';
        successElement = document.getElementById('weight-success');
        errorElement = document.getElementById('weight-error');
        formElement = document.getElementById('weight-form');
        resultElement = document.getElementById('weight-result');
        errorMessageElement = document.getElementById('weight-error-message');
      } else if (form.closest('#tab3')) {
        apiEndpoint = '/api/temperature';
        successElement = document.getElementById('temperature-success');
        errorElement = document.getElementById('temperature-error');
        formElement = document.getElementById('temperature-form');
        resultElement = document.getElementById('temperature-result');
        errorMessageElement = document.getElementById('temperature-error-message');
      } else {
        alert('Unknown form, cannot submit.');
        return;
      }
  
      fetch(apiEndpoint, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      })
        .then(async res => {
          if (!res.ok) {
            const data = await res.json();
              throw new Error(data.error || 'Network response was not ok');
          }
          return res.json();
        })
        .then(json => {
          // Handle success
          if (successElement && formElement && resultElement) {
            // Hide the form and error message
            formElement.classList.add('hidden');
            if (errorElement) errorElement.classList.add('hidden');
            
            // Show success message with result
            resultElement.textContent = json.result || json.message || 'Conversion successful!';
            successElement.classList.remove('hidden');
          } else {
            // Fallback if elements not found
            alert(json.message || 'Submission successful!');
            form.reset();
          }
        })
        .catch(err => {
          // Handle error
          if (errorElement && errorMessageElement) {
            // Hide success message
            if (successElement) successElement.classList.add('hidden');
            
            // Show error message
            errorMessageElement.textContent = err.message || 'Submission failed';
            errorElement.classList.remove('hidden');
          } else {
            // Fallback if elements not found
            alert('Submission failed: ' + err.message);
          }
        });
    }
  
    const forms = document.querySelectorAll('form');
    forms.forEach(form => {
      form.addEventListener('submit', e => {
        e.preventDefault();
        submitForm(form);
      });
    });
    
    // Add event listeners for reset buttons
    const resetButtons = document.querySelectorAll('[id$="-reset"]');
    resetButtons.forEach(button => {
      button.addEventListener('click', () => {
        // Determine which tab we're in
        const tabId = button.closest('[role="tabpanel"]').id;
        
        if (tabId === 'tab1') {
          // Length converter reset
          const successElement = document.getElementById('length-success');
          const errorElement = document.getElementById('length-error');
          const formElement = document.getElementById('length-form');
          
          // Hide success and error messages
          successElement.classList.add('hidden');
          errorElement.classList.add('hidden');
          
          // Show and reset form
          formElement.classList.remove('hidden');
          formElement.reset();
        } else if (tabId === 'tab2') {
          // Weight converter reset
          const successElement = document.getElementById('weight-success');
          const errorElement = document.getElementById('weight-error');
          const formElement = document.getElementById('weight-form');

          // Hide success and error messages
          successElement.classList.add('hidden');
          errorElement.classList.add('hidden');

          // Show and reset form
          formElement.classList.remove('hidden');
          formElement.reset();
        } else if (tabId === 'tab3') {
          // Temperature converter reset
          const successElement = document.getElementById('temperature-success');
          const errorElement = document.getElementById('temperature-error');
          const formElement = document.getElementById('temperature-form');
          // Hide success and error messages
          successElement.classList.add('hidden');
          errorElement.classList.add('hidden');
          // Show and reset form
          formElement.classList.remove('hidden');
          formElement.reset();
        }
      });
    });
});