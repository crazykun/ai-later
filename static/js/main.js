// static/js/main.js
(function() {
    'use strict';

    // ============================================
    // 工具函数
    // ============================================

    // 节流函数 - 用于优化高频事件
    function throttle(func, wait) {
        let timeout;
        let lastRan;
        return function() {
            const context = this;
            const args = arguments;
            if (!lastRan) {
                func.apply(context, args);
                lastRan = Date.now();
            } else {
                clearTimeout(timeout);
                timeout = setTimeout(function() {
                    if ((Date.now() - lastRan) >= wait) {
                        func.apply(context, args);
                        lastRan = Date.now();
                    }
                }, wait - (Date.now() - lastRan));
            }
        };
    }

    // 检查是否为移动端
    function isMobile() {
        return window.innerWidth < 640;
    }

    // ============================================
    // 颜色和首字母生成
    // ============================================

    function generateColorFromName(name) {
        let hash = 0;
        for (let i = 0; i < name.length; i++) {
            hash = name.charCodeAt(i) + ((hash << 5) - hash);
        }
        const h = Math.abs(hash % 360);
        const s = 70;
        const l = 65;
        return hslToHex(h, s, l);
    }

    function hslToHex(h, s, l) {
        l /= 100;
        const a = s * Math.min(l, 1 - l) / 100;
        const f = n => {
            const k = (n + h / 30) % 12;
            const color = l - a * Math.max(Math.min(k - 3, 9 - k, 1), -1);
            return Math.round(255 * color).toString(16).padStart(2, '0');
        };
        return `#${f(0)}${f(8)}${f(4)}`;
    }

    function getInitialsFromName(name) {
        if (!name) return '?';
        const runes = Array.from(name);
        return runes.slice(0, 2).join('');
    }

    function handleImageError(imgElement) {
        const card = imgElement.closest('[data-name]');
        if (!card) return;

        const siteName = card.getAttribute('data-name');
        if (!siteName) return;

        const color = generateColorFromName(siteName);
        const initials = getInitialsFromName(siteName);

        const container = imgElement.parentElement;
        const placeholderDiv = document.createElement('div');
        placeholderDiv.className = imgElement.className;
        placeholderDiv.style.cssText = `
            background-color: ${color};
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-weight: bold;
            font-size: ${imgElement.classList.contains('h-10') && imgElement.classList.contains('w-10') ? '0.875rem' : '1.125rem'};
        `;
        placeholderDiv.textContent = initials;

        container.replaceChild(placeholderDiv, imgElement);
    }

    // ============================================
    // 页面加载和导航栏效果
    // ============================================

    function initPageEffects() {
        // 立即隐藏加载器，移除动画
        const loader = document.getElementById('loader');
        if (loader) {
            loader.classList.add('hidden');
        }

        // 导航栏滚动效果 - 使用节流优化
        const header = document.querySelector('.header-scroll');
        if (header) {
            window.addEventListener('scroll', throttle(function() {
                header.classList.toggle('scrolled', window.scrollY > 50);
            }, 100));
        }
    }

    // ============================================
    // 卡片动画优化
    // ============================================

    function initCardAnimations() {
        // 在页面完全加载后再启用卡片过渡效果
        setTimeout(() => {
            const siteCards = document.querySelectorAll('.grid > div');
            siteCards.forEach(card => {
                // 只添加 hover-card 类，不启用过渡
                card.classList.add('hover-card');
                // 监听第一次 hover 事件后才启用过渡
                card.addEventListener('mouseenter', function() {
                    this.classList.add('transition-enabled');
                }, { once: true });
            });
        }, 100);
    }

    // ============================================
    // 事件委托 - 卡片效果
    // ============================================

    function initCardEffects() {
        const grid = document.querySelector('.grid');
        if (!grid) return;

        // 按钮波纹效果使用事件委托
        grid.addEventListener('click', function(e) {
            const btn = e.target.closest('a[href]');
            if (!btn || !btn.textContent.includes('访问站点')) return;

            const ripple = document.createElement('span');
            const rect = btn.getBoundingClientRect();
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
                animation: ripple 0.5s ease-out forwards;
                pointer-events: none;
            `;

            btn.style.position = 'relative';
            btn.style.overflow = 'hidden';
            btn.appendChild(ripple);

            setTimeout(() => ripple.remove(), 500);
        });
    }

    // ============================================
    // 移动端菜单和搜索
    // ============================================

    function initMobileMenu() {
        const searchBtn = document.getElementById('searchBtn');
        const menuBtn = document.getElementById('menuBtn');
        const searchBox = document.getElementById('searchBox');
        const navMenu = document.getElementById('navMenu');
        const menuIcon = document.getElementById('menuIcon');
        const closeIcon = document.getElementById('closeIcon');

        function toggleMenuIcon(isOpen) {
            if (menuIcon && closeIcon) {
                menuIcon.classList.toggle('hidden', isOpen);
                closeIcon.classList.toggle('hidden', !isOpen);
            }
        }

        // 搜索框切换
        if (searchBtn && searchBox) {
            searchBtn.addEventListener('click', function(e) {
                e.stopPropagation();
                if (!isMobile()) return;

                const isOpen = !searchBox.classList.toggle('hidden');
                if (isOpen) {
                    const input = searchBox.querySelector('input[name="q"]');
                    if (input) input.focus();
                    if (navMenu) {
                        navMenu.classList.add('hidden');
                        toggleMenuIcon(false);
                    }
                }
            });
        }

        // 菜单切换
        if (menuBtn && navMenu) {
            menuBtn.addEventListener('click', function(e) {
                e.stopPropagation();
                if (!isMobile()) return;

                const isOpen = !navMenu.classList.contains('hidden');
                if (isOpen) {
                    navMenu.classList.add('hidden');
                    toggleMenuIcon(false);
                } else {
                    navMenu.classList.remove('hidden');
                    toggleMenuIcon(true);
                    if (searchBox) searchBox.classList.add('hidden');
                }
            });
        }

        // 点击外部关闭
        document.addEventListener('click', function(e) {
            if (!isMobile()) return;

            if (navMenu && !navMenu.classList.contains('hidden')) {
                if (!navMenu.contains(e.target) && !menuBtn?.contains(e.target)) {
                    navMenu.classList.add('hidden');
                    toggleMenuIcon(false);
                }
            }

            if (searchBox && !searchBox.classList.contains('hidden')) {
                if (!searchBox.contains(e.target) && !searchBtn?.contains(e.target)) {
                    searchBox.classList.add('hidden');
                }
            }
        });

        // 导航链接点击关闭
        if (navMenu) {
            navMenu.querySelectorAll('a').forEach(link => {
                link.addEventListener('click', function() {
                    if (isMobile()) {
                        navMenu.classList.add('hidden');
                        toggleMenuIcon(false);
                        if (searchBox) searchBox.classList.add('hidden');
                    }
                });
            });
        }

        // 窗口大小改变时重置 - 使用节流
        window.addEventListener('resize', throttle(function() {
            if (!isMobile()) {
                if (navMenu) {
                    navMenu.classList.add('hidden');
                    toggleMenuIcon(false);
                }
                if (searchBox) searchBox.classList.add('hidden');
            }
        }, 200));

        // ESC 键关闭
        document.addEventListener('keydown', function(e) {
            if (e.key === 'Escape' && isMobile()) {
                if (navMenu && !navMenu.classList.contains('hidden')) {
                    navMenu.classList.add('hidden');
                    toggleMenuIcon(false);
                }
                if (searchBox && !searchBox.classList.contains('hidden')) {
                    searchBox.classList.add('hidden');
                }
            }
        });
    }

    // ============================================
    // 返回顶部按钮
    // ============================================

    function initScrollTopButton() {
        let scrollTopBtn = null;

        function createButton() {
            if (scrollTopBtn) return;

            scrollTopBtn = document.createElement('button');
            scrollTopBtn.innerHTML = `
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 10l7-7m0 0l7 7m-7-7v18"></path>
                </svg>
            `;
            scrollTopBtn.className = 'fixed bottom-6 right-6 bg-blue-500 text-white p-3 rounded-full shadow-lg hover:bg-blue-600 transition-opacity duration-200 z-50 opacity-0 pointer-events-none';
            scrollTopBtn.onclick = () => {
                window.scrollTo({ top: 0, behavior: 'smooth' });
            };
            document.body.appendChild(scrollTopBtn);
        }

        // 优化：使用节流减少滚动事件处理频率
        window.addEventListener('scroll', throttle(function() {
            if (!scrollTopBtn) createButton();

            const show = window.scrollY > 300;
            scrollTopBtn.classList.toggle('opacity-0', !show);
            scrollTopBtn.classList.toggle('opacity-100', show);
            scrollTopBtn.classList.toggle('pointer-events-none', !show);
        }, 100));
    }

    // ============================================
    // 其他功能
    // ============================================

    function initOtherFeatures() {
        // 分类筛选
        const categorySelect = document.querySelector('select[name="category"]');
        if (categorySelect) {
            categorySelect.addEventListener('change', function() {
                this.closest('form').submit();
            });
        }

        // 左侧导航动画
        const sidebarNav = document.querySelector('.sidebar-nav');
        if (sidebarNav) {
            sidebarNav.classList.add('sidebar-nav');
        }

        // 图片错误处理
        const siteImages = document.querySelectorAll('.grid > div img[src]');
        siteImages.forEach(img => {
            img.onerror = function() {
                handleImageError(this);
            };
        });
    }

    // ============================================
    // 初始化
    // ============================================

    // 添加必要的 CSS 动画
    const style = document.createElement('style');
    style.textContent = `
        @keyframes ripple {
            to {
                transform: scale(3);
                opacity: 0;
            }
        }
    `;
    document.head.appendChild(style);

    // DOM 加载完成后初始化
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', initAll);
    } else {
        initAll();
    }

    function initAll() {
        initPageEffects();
        initCardAnimations();
        initCardEffects();
        initMobileMenu();
        initScrollTopButton();
        initOtherFeatures();
    }

})();
