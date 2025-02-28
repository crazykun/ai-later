// static/js/main.js
document.addEventListener('DOMContentLoaded', function() {
    // Add fade-in animation to site cards
    const siteCards = document.querySelectorAll('.grid > div');
    siteCards.forEach((card, index) => {
        card.classList.add('fade-in');
        card.style.animationDelay = `${index * 0.1}s`;
    });

    // Add hover effect class
    siteCards.forEach(card => {
        card.classList.add('hover-card');
    });

    // Handle category filter changes
    const categorySelect = document.querySelector('select[name="category"]');
    if (categorySelect) {
        categorySelect.addEventListener('change', function() {
            this.closest('form').submit();
        });
    }
});