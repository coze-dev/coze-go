// ==UserScript==
// @name         移除 Feed Good 组件
// @namespace    http://tampermonkey.net/
// @version      1.1
// @description  移除页面中的 Feed Good 调查组件
// @author       You
// @match        *://data.bytedance.net/*
// @grant        none
// ==/UserScript==

(function() {
    'use strict';

    const selectors = [
        '[data-feelgood-task-id]'
    ];

    function removeFeedGood() {
        selectors.forEach(selector => {
            const elements = document.querySelectorAll(selector);
            elements.forEach(element => {
                // 检查是否包含特定的图标或特征
                const hasFeelGoodIcon = element.querySelector('img[src*="feelgood-open/fg-icon-chat.svg"]');
                const hasFeelGoodTaskId = element.getAttribute('data-feelgood-task-id');
                const hasAthenaClass = element.classList.contains('athena-survey-entry') || 
                                      element.classList.contains('athena-survey-widget');
                
                if (hasFeelGoodIcon || hasFeelGoodTaskId || hasAthenaClass) {
                    console.log('发现 Feed Good 组件，正在移除...');
                    element.remove();
                }
            });
        });
    }

    // 初始执行
    removeFeedGood();

    // 监听 DOM 变化，处理动态加载的组件
    const observer = new MutationObserver(function(mutations) {
        mutations.forEach(function(mutation) {
            if (mutation.type === 'childList') {
                mutation.addedNodes.forEach(function(node) {
                    if (node.nodeType === Node.ELEMENT_NODE) {
                        selectors.forEach(selector => {
                            const feedGoodElements = node.querySelectorAll ? 
                                node.querySelectorAll(selector) : [];
                            
                            feedGoodElements.forEach(element => {
                                const hasFeelGoodIcon = element.querySelector('img[src*="feelgood-open/fg-icon-chat.svg"]');
                                const hasFeelGoodTaskId = element.getAttribute('data-feelgood-task-id');
                                const hasAthenaClass = element.classList.contains('athena-survey-entry') || 
                                                      element.classList.contains('athena-survey-widget');
                                
                                if (hasFeelGoodIcon || hasFeelGoodTaskId || hasAthenaClass) {
                                    console.log('发现动态加载的 Feed Good 组件，正在移除...');
                                    element.remove();
                                }
                            });
                        });

                        // 检查节点本身是否是 Feed Good 组件
                        if (node.classList) {
                            const hasFeelGoodIcon = node.querySelector('img[src*="feelgood-open/fg-icon-chat.svg"]');
                            const hasFeelGoodTaskId = node.getAttribute('data-feelgood-task-id');
                            const hasAthenaClass = node.classList.contains('athena-survey-entry') || 
                                                  node.classList.contains('athena-survey-widget');
                            
                            if (hasFeelGoodIcon || hasFeelGoodTaskId || hasAthenaClass) {
                                console.log('发现动态加载的 Feed Good 组件，正在移除...');
                                node.remove();
                            }
                        }
                    }
                });
            }
        });
    });

    // 开始观察 DOM 变化
    observer.observe(document.body, {
        childList: true,
        subtree: true
    });

    // 定期检查，以防某些组件没有被 MutationObserver 捕获
    setInterval(removeFeedGood, 1000);

    console.log('Feed Good 组件移除脚本已加载 (v1.1)');
})(); 