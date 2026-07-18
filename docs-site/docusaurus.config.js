// @ts-check

/** @type {import('@docusaurus/types').Config} */
const config = {
	title: "HyperNexus",
	tagline: "AI Control Plane with Persistent Memory & 26,000+ MCP Tools",
	favicon: "img/favicon.ico",
	url: "https://docs.hypernexus.site",
	baseUrl: "/",
	organizationName: "MDMAtk",
	projectName: "HyperNexus",
	onBrokenLinks: "throw",
	onBrokenMarkdownLinks: "warn",

	presets: [
		[
			"classic",
			/** @type {import('@docusaurus/preset-classic').Options} */
			({
				docs: {
					sidebarPath: "./sidebars.js",
					editUrl:
						"https://github.com/MDMAtk/HyperNexus/tree/main/docs-site/",
				},
				blog: {
					showReadingTime: true,
					editUrl:
						"https://github.com/MDMAtk/HyperNexus/tree/main/docs-site/",
				},
				theme: {
					customCss: "./src/css/custom.css",
				},
			}),
		],
	],

	themeConfig:
		/** @type {import('@docusaurus/preset-classic').ThemeConfig} */
		({
			colorMode: {
				defaultMode: "dark",
				respectPrefersColorScheme: true,
			},
			navbar: {
				title: "HyperNexus",
				logo: {
					alt: "HyperNexus Logo",
					src: "img/logo.svg",
				},
				items: [
					{
						type: "doc",
						docId: "intro",
						position: "left",
						label: "Docs",
					},
					{ to: "/blog", label: "Blog", position: "left" },
					{
						href: "https://github.com/MDMAtk/HyperNexus",
						label: "GitHub",
						position: "right",
					},
					{
						href: "https://discord.gg/Hj9P3GbVxR",
						label: "Discord",
						position: "right",
					},
				],
			},
			footer: {
				style: "dark",
				links: [
					{
						title: "Docs",
						items: [
							{ label: "Getting Started", to: "/docs/intro" },
							{ label: "API Reference", to: "/docs/api" },
						],
					},
					{
						title: "Community",
						items: [
							{ label: "Discord", href: "https://discord.gg/Hj9P3GbVxR" },
							{
								label: "GitHub",
								href: "https://github.com/MDMAtk/HyperNexus",
							},
						],
					},
					{
						title: "More",
						items: [
							{ label: "Blog", to: "/blog" },
							{ label: "Website", href: "https://hypernexus.site" },
						],
					},
				],
				copyright: `Copyright © ${new Date().getFullYear()} HyperNexus. Built with Docusaurus.`,
			},
			prism: {
				theme: require("prism-react-renderer").themes.github,
				darkTheme: require("prism-react-renderer").themes.dracula,
			},
		}),
};

module.exports = config;
