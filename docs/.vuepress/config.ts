import {defineConfig} from 'vuepress/config'

export default defineConfig(ctx => ({
    theme: '@vuepress/theme-vue',
    head: [
        ['link', {rel: 'icon', href: `/logo.png`}],
        ['link', {rel: 'manifest', href: '/manifest.json'}],
        ['meta', {name: 'theme-color', content: '#3eaf7c'}],
        ['meta', {name: 'apple-mobile-web-app-capable', content: 'yes'}],
        [
            'meta',
            {name: 'apple-mobile-web-app-status-bar-style', content: 'black'}
        ],
        [
            'link',
            {rel: 'apple-touch-icon', href: `/icons/apple-touch-icon-152x152.png`}
        ],
        [
            'link',
            {
                rel: 'mask-icon',
                href: '/icons/safari-pinned-tab.svg',
                color: '#3eaf7c'
            }
        ],
        [
            'meta',
            {
                name: 'msapplication-TileImage',
                content: '/icons/msapplication-icon-144x144.png'
            }
        ],
        ['meta', {name: 'msapplication-TileColor', content: '#000000'}]
    ],
    locales: {
        '/': {
            lang: 'en-US',
            title: 'Weidows/Golang',
            description: 'Vue-powered Static Site Generator'
        },
        '/zh/': {
            lang: 'zh-CN',
            title: 'Weidows/Golang',
            description: 'Vue 驱动的静态网站生成器'
        }
    },
    themeConfig: {
        repo: 'Weidows/Golang',
        editLinks: true,
        docsDir: 'docs',
        // #697 Provided by the official algolia team.
        algolia: ctx.isProd
            ? {
                apiKey: '3a539aab83105f01761a137c61004d85',
                indexName: 'vuepress',
                algoliaOptions: {
                    facetFilters: ['tags:v1']
                }
            }
            : null,
        smoothScroll: true,
        locales: {
            '/': {
                label: 'English',
                selectText: 'Languages',
                ariaLabel: 'Select language',
                editLinkText: 'Edit this page on GitHub',
                lastUpdated: 'Last Updated',
            },
            '/zh/': {
                label: '简体中文',
                selectText: '选择语言',
                ariaLabel: '选择语言',
                editLinkText: '在 GitHub 上编辑此页',
                lastUpdated: '上次更新',
            }
        },
        nav: [
            {
                text: "首页",
                link: "/",
            },
        ],
        sidebar: [
            {
                title: "/",
                path: "/",
                collapsable: true,
                sidebarDepth: 2,
                children: [
                    {
                        title: "CMD",
                        path: "/cmd/",
                        children: [
                            "/cmd/dsg",
                            "/cmd/gmm",
                            "/cmd/jpu",
                        ],
                    },
                    {
                        title: "Utils",
                        path: "/utils/",
                        children: [
                        ],
                    },
                ],
            },
        ],
    },
    plugins: [
        ['@vuepress/back-to-top', true],
        [
            '@vuepress/pwa',
            {
                serviceWorker: true,
                updatePopup: true
            }
        ],
        ['@vuepress/medium-zoom', true],
        [
            '@vuepress/google-analytics',
            {
                ga: 'UA-128189152-1'
            }
        ],
        [
            'vuepress-plugin-container',
            {
                type: 'vue',
                before: '<pre class="vue-container"><code>',
                after: '</code></pre>'
            }
        ],
        [
            'vuepress-plugin-container',
            {
                type: 'upgrade',
                before: info => `<UpgradePath title="${info}">`,
                after: '</UpgradePath>'
            }
        ],
        ['vuepress-plugin-flowchart']
    ],
    evergreen: !ctx.isProd
}))