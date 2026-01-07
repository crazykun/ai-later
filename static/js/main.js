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

    // 获取移动端元素
    const searchBtn = document.getElementById('searchBtn');
    const menuBtn = document.getElementById('menuBtn');
    const searchBox = document.getElementById('searchBox');
    const navMenu = document.getElementById('navMenu');
    const menuIcon = document.getElementById('menuIcon');
    const closeIcon = document.getElementById('closeIcon');
    
    // 检查是否为移动端
    function isMobile() {
        return window.innerWidth < 640; // sm breakpoint
    }
    
    // 移动端搜索框切换
    if (searchBtn && searchBox) {
        searchBtn.addEventListener('click', function(e) {
            e.stopPropagation();
            if (!isMobile()) return; // 只在移动端工作
            
            searchBox.classList.toggle('hidden');
            // 关闭导航菜单
            if (navMenu && !navMenu.classList.contains('hidden')) {
                navMenu.classList.add('hidden');
                toggleMenuIcon(false);
            }
            
            if (!searchBox.classList.contains('hidden')) {
                const searchInput = searchBox.querySelector('input[name="q"]');
                if (searchInput) {
                    searchInput.focus();
                }
            }
        });
    }

    // 移动端导航菜单切换
    if (menuBtn && navMenu) {
        menuBtn.addEventListener('click', function(e) {
            e.stopPropagation();
            if (!isMobile()) return; // 只在移动端工作
            
            const isMenuOpen = !navMenu.classList.contains('hidden');
            
            if (isMenuOpen) {
                navMenu.classList.add('hidden');
                toggleMenuIcon(false);
            } else {
                navMenu.classList.remove('hidden');
                toggleMenuIcon(true);
                // 关闭搜索框
                if (searchBox && !searchBox.classList.contains('hidden')) {
                    searchBox.classList.add('hidden');
                }
            }
        });
    }

    // 切换菜单图标
    function toggleMenuIcon(isOpen) {
        if (menuIcon && closeIcon) {
            if (isOpen) {
                menuIcon.classList.add('hidden');
                closeIcon.classList.remove('hidden');
            } else {
                menuIcon.classList.remove('hidden');
                closeIcon.classList.add('hidden');
            }
        }
    }

    // 点击页面其他地方关闭移动端菜单和搜索框
    document.addEventListener('click', function(event) {
        if (!isMobile()) return; // 只在移动端工作
        
        const isClickInsideMenu = navMenu && navMenu.contains(event.target);
        const isClickInsideMenuBtn = menuBtn && menuBtn.contains(event.target);
        const isClickInsideSearch = searchBox && searchBox.contains(event.target);
        const isClickInsideSearchBtn = searchBtn && searchBtn.contains(event.target);
        
        if (!isClickInsideMenu && !isClickInsideMenuBtn && navMenu && !navMenu.classList.contains('hidden')) {
            navMenu.classList.add('hidden');
            toggleMenuIcon(false);
        }
        
        if (!isClickInsideSearch && !isClickInsideSearchBtn && searchBox && !searchBox.classList.contains('hidden')) {
            searchBox.classList.add('hidden');
        }
    });

    // 移动端导航链接点击后关闭菜单
    const navLinks = navMenu ? navMenu.querySelectorAll('a') : [];
    navLinks.forEach(link => {
        link.addEventListener('click', function() {
            if (!isMobile()) return; // 只在移动端工作
            
            if (navMenu) {
                navMenu.classList.add('hidden');
                toggleMenuIcon(false);
            }
            if (searchBox) {
                searchBox.classList.add('hidden');
            }
        });
    });

    // 窗口大小改变时重置移动端菜单状态
    window.addEventListener('resize', function() {
        if (!isMobile()) {
            // 切换到PC端时，隐藏移动端菜单
            if (navMenu) {
                navMenu.classList.add('hidden');
                toggleMenuIcon(false);
            }
            if (searchBox) {
                searchBox.classList.add('hidden');
            }
        }
    });

    // 左侧导航动画
    const sidebarNav = document.querySelector('.sidebar-nav');
    if (sidebarNav) {
        sidebarNav.classList.add('sidebar-nav');
    }

    // 为卡片添加更好的悬停效果和按钮动画
    siteCards.forEach(card => {
        card.addEventListener('mouseenter', function() {
            this.style.transform = 'translateY(-8px) scale(1.02)';
            
            // 为访问站点按钮添加特殊效果
            const visitBtn = this.querySelector('a[href]');
            if (visitBtn) {
                visitBtn.classList.add('visit-btn');
            }
        });
        
        card.addEventListener('mouseleave', function() {
            this.style.transform = 'translateY(0) scale(1)';
        });
        
        // 为按钮添加点击波纹效果
        const visitBtn = card.querySelector('a[href]');
        if (visitBtn) {
            visitBtn.addEventListener('click', function(e) {
                // 创建波纹效果
                const ripple = document.createElement('span');
                const rect = this.getBoundingClientRect();
                const size = Math.max(rect.width, rect.height);
                const x = e.clientX - rect.left - size / 2;
                const y = e.clientY - rect.top - size / 2;
                
                ripple.style.cssText = `
                    position: absolute;
                    width: ${size}px;
                    height: ${size}px;
                    left: ${x}px;
                    top: ${y}px;
                    background: rgba(255, 255, 255, 0.3);
                    border-radius: 50%;
                    transform: scale(0);
                    animation: ripple 0.6s linear;
                    pointer-events: none;
                `;
                
                this.style.position = 'relative';
                this.style.overflow = 'hidden';
                this.appendChild(ripple);
                
                setTimeout(() => {
                    ripple.remove();
                }, 600);
            });
        }
    });

    // 添加波纹动画CSS
    const style = document.createElement('style');
    style.textContent = `
        @keyframes ripple {
            to {
                transform: scale(4);
                opacity: 0;
            }
        }
        
        .visit-btn {
            position: relative;
            overflow: hidden;
        }
        
        .card-enter {
            animation: cardEnter 0.6s ease-out forwards;
        }
        
        @keyframes cardEnter {
            from {
                opacity: 0;
                transform: translateY(30px) scale(0.9);
            }
            to {
                opacity: 1;
                transform: translateY(0) scale(1);
            }
        }
    `;
    document.head.appendChild(style);

    // 为新加载的卡片添加进入动画
    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.classList.add('card-enter');
                observer.unobserve(entry.target);
            }
        });
    }, {
        threshold: 0.1,
        rootMargin: '0px 0px -50px 0px'
    });

    siteCards.forEach(card => {
        observer.observe(card);
    });

    // 平滑滚动到顶部
    function scrollToTop() {
        window.scrollTo({
            top: 0,
            behavior: 'smooth'
        });
    }

    // 添加返回顶部按钮（如果页面较长）
    let scrollTopBtn = null;
    function createScrollTopButton() {
        if (!scrollTopBtn) {
            scrollTopBtn = document.createElement('button');
            scrollTopBtn.innerHTML = `
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 10l7-7m0 0l7 7m-7-7v18"></path>
                </svg>
            `;
            scrollTopBtn.className = 'fixed bottom-6 right-6 bg-blue-500 text-white p-3 rounded-full shadow-lg hover:bg-blue-600 transition-all duration-200 z-50 opacity-0 pointer-events-none';
            scrollTopBtn.onclick = scrollToTop;
            document.body.appendChild(scrollTopBtn);
        }
    }

    // 监听滚动事件
    window.addEventListener('scroll', function() {
        if (!scrollTopBtn) createScrollTopButton();
        
        if (window.scrollY > 300) {
            scrollTopBtn.classList.remove('opacity-0', 'pointer-events-none');
            scrollTopBtn.classList.add('opacity-100');
        } else {
            scrollTopBtn.classList.add('opacity-0', 'pointer-events-none');
            scrollTopBtn.classList.remove('opacity-100');
        }
    });

    // ESC键关闭移动端菜单
    document.addEventListener('keydown', function(event) {
        if (event.key === 'Escape' && isMobile()) {
            if (navMenu && !navMenu.classList.contains('hidden')) {
                navMenu.classList.add('hidden');
                toggleMenuIcon(false);
            }
            if (searchBox && !searchBox.classList.contains('hidden')) {
                searchBox.classList.add('hidden');
            }
        }
    });
});