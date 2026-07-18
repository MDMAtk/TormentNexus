import { Info } from "lucide-react";

export default function MockupPage() {
    return (
        <div className="min-h-screen bg-black text-zinc-100 p-6 flex justify-center">
            {/* 1200px max width container as requested */}
            <div className="w-full max-w-[1200px] flex flex-col gap-6">

                {/* Header Section */}
                <div className="flex items-center justify-between border-b border-zinc-800 pb-4">
                    <h1 className="text-2xl font-bold">TormentNexus Dashboard</h1>
                    <div className="text-sm text-zinc-500">System Live</div>
                </div>

                {/* Top Tier: Two-Column Grid */}
                <div className="grid grid-cols-1 xl:grid-cols-2 gap-6">
                    {/* Left Column: Memory/Context + Workflows */}
                    <div className="flex flex-col gap-6">
                        <section className="bg-zinc-900/50 border border-zinc-800 rounded-lg p-5">
                            <div className="flex items-center gap-2 mb-4">
                                <h2 className="text-lg font-semibold">Memory & Context</h2>
                                <span title="View active session memories, L1/L2 storage, and RAG context" className="cursor-help">💡</span>
                            </div>
                            <div className="h-48 border border-dashed border-zinc-700 flex items-center justify-center text-zinc-500 rounded">
                                Memory Widget Placeholder
                            </div>
                        </section>

                        <section className="bg-zinc-900/50 border border-zinc-800 rounded-lg p-5">
                            <div className="flex items-center gap-2 mb-4">
                                <h2 className="text-lg font-semibold">Workflows & Swarm</h2>
                                <span title="Monitor active Swarm workers and orchestrate multi-agent workflows" className="cursor-help">💡</span>
                            </div>
                            <div className="h-48 border border-dashed border-zinc-700 flex items-center justify-center text-zinc-500 rounded">
                                Workflows Widget Placeholder
                            </div>
                        </section>
                    </div>

                    {/* Right Column: MCP/Integrations + Tool Consoles */}
                    <div className="flex flex-col gap-6">
                        <section className="bg-zinc-900/50 border border-zinc-800 rounded-lg p-5">
                            <div className="flex items-center gap-2 mb-4">
                                <h2 className="text-lg font-semibold">MCP & Integrations</h2>
                                <span title="Manage Model Context Protocol servers, connections, and external API integrations" className="cursor-help">💡</span>
                            </div>
                            <div className="h-48 border border-dashed border-zinc-700 flex items-center justify-center text-zinc-500 rounded">
                                MCP Registry Placeholder
                            </div>
                        </section>

                        <section className="bg-zinc-900/50 border border-zinc-800 rounded-lg p-5">
                            <div className="flex items-center gap-2 mb-4">
                                <h2 className="text-lg font-semibold">Tool Consoles</h2>
                                <span title="Direct access to command surfaces, terminal logs, and system console" className="cursor-help">💡</span>
                            </div>
                            <div className="h-48 border border-dashed border-zinc-700 flex items-center justify-center text-zinc-500 rounded">
                                Console / Terminal Placeholder
                            </div>
                        </section>
                    </div>
                </div>

                {/* Expandable Cards for Lower Priority Modules */}
                <div className="mt-8 space-y-4">
                    <details className="group bg-zinc-900/30 border border-zinc-800 rounded-lg overflow-hidden">
                        <summary className="flex items-center justify-between p-4 cursor-pointer hover:bg-zinc-800/50 transition-colors">
                            <div className="flex items-center gap-2">
                                <span className="font-semibold text-zinc-300 group-open:text-white">Code Sandbox</span>
                                <span title="Isolated execution environment for arbitrary code blocks" className="cursor-help">💡</span>
                            </div>
                            <span className="text-zinc-500 group-open:rotate-180 transition-transform">▼</span>
                        </summary>
                        <div className="p-4 border-t border-zinc-800 text-zinc-400">
                            Code Sandbox UI goes here...
                        </div>
                    </details>

                    <details className="group bg-zinc-900/30 border border-zinc-800 rounded-lg overflow-hidden">
                        <summary className="flex items-center justify-between p-4 cursor-pointer hover:bg-zinc-800/50 transition-colors">
                            <div className="flex items-center gap-2">
                                <span className="font-semibold text-zinc-300 group-open:text-white">System Health & Settings</span>
                                <span title="Resource metrics, commercial settings, and global configuration" className="cursor-help">💡</span>
                            </div>
                            <span className="text-zinc-500 group-open:rotate-180 transition-transform">▼</span>
                        </summary>
                        <div className="p-4 border-t border-zinc-800 text-zinc-400">
                            Settings and Health Modules go here...
                        </div>
                    </details>
                </div>

            </div>
        </div>
    );
}
