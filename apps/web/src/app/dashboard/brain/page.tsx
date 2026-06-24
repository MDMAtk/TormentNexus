'use client';

import React, { useState } from 'react';
import { Card, CardHeader, CardTitle, CardContent, Button, Input, Tabs, TabsContent, TabsList, TabsTrigger, KnowledgeGraph } from '@tormentnexus/ui';
import { Loader2, RefreshCw, Download, Play, CheckCircle, Search, Sparkles, Bot, Globe } from "lucide-react";
import { trpc } from '@/utils/trpc';
import { toast } from "sonner";

type ExpertTrpc = {
    expert: {
        research: { useMutation: () => any };
        code: { useMutation: () => any };
    };
};

export default function BrainPage() {
    const trpcWithExpert = trpc as unknown as typeof trpc & ExpertTrpc;
    const graphQuery = trpc.graph.getSymbolsGraph.useQuery();
    const { nodes: rawNodes = [], links: rawLinks = [] } = graphQuery.data || {};

    const nodes = rawNodes.map((node: any) => ({
        id: node.id,
        label: node.name,
        type: 'concept' as const,
        val: node.val || 1
    }));

    const links = rawLinks.map((link: any) => ({
        source: link.source,
        target: link.target,
        value: 1
    }));

    // Ingest state
    const [ingestUrl, setIngestUrl] = useState("");
    const [ingestLog, setIngestLog] = useState("");
    const ingestMutation = trpc.knowledge.ingest.useMutation();
    const resourcesQuery = trpc.knowledge.getResources.useQuery();
    const resources = resourcesQuery.data || { categories: [] };

    // Research state
    const [researchQuery, setResearchQuery] = useState("");
    const [researchDepth, setResearchDepth] = useState(2);
    const researchMutation = trpcWithExpert.expert.research.useMutation();

    // Coder state
    const [coderTask, setCoderTask] = useState("");
    const coderMutation = trpcWithExpert.expert.code.useMutation({
        onSuccess: () => {
            toast.success("Coder task started");
            setCoderTask("");
        },
        onError: (err) => toast.error("Coder task failed: " + err.message)
    });

    const handleIngest = async () => {
        if (!ingestUrl) return;
        setIngestLog(`Ingesting: ${ingestUrl}...`);
        try {
            const result = await ingestMutation.mutateAsync({ url: ingestUrl });
            setIngestLog(`Success: ${result}`);
            setIngestUrl("");
            resourcesQuery.refetch();
        } catch (e: any) {
            setIngestLog(`Error: ${e.message}`);
        }
    };

    const handleResearch = () => {
        if (!researchQuery) return;
        researchMutation.mutate({ query: researchQuery, depth: researchDepth, breadth: 3 });
    };

    const handleCode = () => {
        if (!coderTask) return;
        coderMutation.mutate({ task: coderTask });
    };

    return (
        <div className="h-full w-full p-6 flex flex-col space-y-6">
            <header className="flex justify-between items-center border-b border-zinc-800 pb-4">
                <div>
                    <h1 className="text-3xl font-bold bg-gradient-to-r from-purple-500 to-pink-500 bg-clip-text text-transparent">
                        TormentNexus Brain &amp; Knowledge
                    </h1>
                    <p className="text-zinc-500 dark:text-zinc-400">
                        Visualizing semantic relationships, ingesting external doc sources, and running expert agents.
                    </p>
                </div>
            </header>

            <Tabs defaultValue="graph" className="w-full flex-1 flex flex-col min-h-0">
                <TabsList className="grid grid-cols-3 max-w-[500px] mb-4 bg-zinc-900 border border-zinc-800 p-1 rounded-lg">
                    <TabsTrigger value="graph" className="text-sm font-medium py-1.5 rounded-md transition-all">Cognitive Graph</TabsTrigger>
                    <TabsTrigger value="ingest" className="text-sm font-medium py-1.5 rounded-md transition-all">URL Ingestion</TabsTrigger>
                    <TabsTrigger value="agents" className="text-sm font-medium py-1.5 rounded-md transition-all">Expert Agents</TabsTrigger>
                </TabsList>

                <TabsContent value="graph" className="flex-1 flex flex-col min-h-0 outline-none relative">
                    <div className="absolute inset-0 bg-zinc-950 border border-zinc-850 rounded-2xl overflow-hidden shadow-inner">
                        <KnowledgeGraph
                            nodes={nodes}
                            links={links}
                            loading={graphQuery.isLoading}
                        />
                    </div>
                </TabsContent>

                <TabsContent value="ingest" className="flex-1 flex flex-col min-h-0 outline-none">
                    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 flex-1 min-h-0 overflow-y-auto">
                        <Card className="bg-zinc-900 border-zinc-800 p-6 flex flex-col gap-4">
                            <CardHeader className="p-0">
                                <CardTitle className="text-lg font-bold text-green-400 flex items-center gap-2">
                                    <Globe className="h-5 w-5" /> Ingest Knowledge Source
                                </CardTitle>
                            </CardHeader>
                            <CardContent className="p-0 flex flex-col gap-4">
                                <p className="text-sm text-zinc-400">
                                    Parse documentation, API references, or articles directly into the TormentNexus context engine.
                                </p>
                                <div className="flex gap-2">
                                    <Input
                                        type="text"
                                        className="flex-1 bg-black border-zinc-800 text-white outline-none placeholder:text-zinc-650"
                                        placeholder="Enter URL (e.g. https://mcp.dev/docs)"
                                        value={ingestUrl}
                                        onChange={(e) => setIngestUrl(e.target.value)}
                                        onKeyDown={(e) => e.key === 'Enter' && handleIngest()}
                                    />
                                    <Button
                                        onClick={handleIngest}
                                        disabled={ingestMutation.isPending}
                                        className="bg-green-600 hover:bg-green-500 text-white font-semibold"
                                    >
                                        {ingestMutation.isPending ? 'Ingesting...' : 'Ingest'}
                                    </Button>
                                </div>
                                {ingestLog && (
                                    <div className="bg-black p-3 rounded-lg text-xs font-mono text-zinc-300 break-all border border-zinc-800 max-h-[250px] overflow-y-auto leading-relaxed">
                                        {ingestLog}
                                    </div>
                                )}
                            </CardContent>
                        </Card>

                        <Card className="bg-zinc-900 border-zinc-800 p-6 flex flex-col gap-4">
                            <CardHeader className="p-0 flex justify-between flex-row items-center">
                                <CardTitle className="text-lg font-bold text-zinc-300">
                                    Ingested Resource Index
                                </CardTitle>
                                <span className="text-xs text-zinc-500">
                                    Last Sync: {resources.lastUpdated ? new Date(resources.lastUpdated).toLocaleTimeString() : 'Never'}
                                </span>
                            </CardHeader>
                            <CardContent className="p-0 flex-1 overflow-y-auto">
                                {resourcesQuery.isLoading ? (
                                    <div className="flex justify-center p-8"><Loader2 className="animate-spin text-zinc-500" /></div>
                                ) : resources.categories?.length === 0 ? (
                                    <p className="text-zinc-500 italic text-center p-8">No external doc sources ingested yet.</p>
                                ) : (
                                    <div className="space-y-4">
                                        {resources.categories?.map((cat: any) => (
                                            <div key={cat.name} className="border-l-2 border-zinc-800 pl-4 py-1">
                                                <h4 className="font-semibold text-zinc-200">{cat.name}</h4>
                                                <ul className="mt-1 space-y-1 text-xs text-zinc-400">
                                                    {cat.items?.map((item: any) => (
                                                        <li key={item.url} className="truncate">
                                                            <a href={item.url} target="_blank" rel="noopener noreferrer" className="hover:text-blue-400 hover:underline">
                                                                {item.title || item.url}
                                                            </a>
                                                        </li>
                                                    ))}
                                                </ul>
                                            </div>
                                        ))}
                                    </div>
                                )}
                            </CardContent>
                        </Card>
                    </div>
                </TabsContent>

                <TabsContent value="agents" className="flex-1 flex flex-col min-h-0 outline-none">
                    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 flex-1 min-h-0 overflow-y-auto">
                        {/* Deep Research Section */}
                        <Card className="bg-zinc-900 border-zinc-800 p-6 flex flex-col gap-4">
                            <CardHeader className="p-0">
                                <CardTitle className="text-lg font-bold text-blue-400 flex items-center gap-2">
                                    <Bot className="h-5 w-5" /> Deep Research Agent
                                </CardTitle>
                            </CardHeader>
                            <CardContent className="p-0 flex flex-col gap-4 flex-1">
                                <div className="flex gap-2">
                                    <Input
                                        type="text"
                                        className="flex-1 bg-black border-zinc-800 text-white outline-none placeholder:text-zinc-650"
                                        placeholder="Research query..."
                                        value={researchQuery}
                                        onChange={(e) => setResearchQuery(e.target.value)}
                                        onKeyDown={(e) => e.key === 'Enter' && handleResearch()}
                                    />
                                    <div className="flex items-center gap-2 bg-black border border-zinc-800 rounded-lg px-3">
                                        <span className="text-[10px] text-zinc-500 font-semibold tracking-wider">DEPTH</span>
                                        <input
                                            type="number"
                                            min="1" max="5"
                                            aria-label="Research depth"
                                            value={researchDepth}
                                            onChange={(e) => setResearchDepth(parseInt(e.target.value))}
                                            className="bg-transparent w-8 text-center outline-none text-white font-semibold"
                                        />
                                    </div>
                                    <Button
                                        onClick={handleResearch}
                                        disabled={researchMutation.isPending}
                                        className="bg-blue-600 hover:bg-blue-500 text-white font-semibold"
                                    >
                                        Research
                                    </Button>
                                </div>

                                <div className="flex-1 overflow-y-auto">
                                    {researchMutation.isPending && (
                                        <div className="flex justify-center p-8"><Loader2 className="animate-spin text-zinc-500" /></div>
                                    )}
                                    {researchMutation.data && (
                                        <div className="bg-black p-4 border border-zinc-800 rounded-lg text-sm text-zinc-300 whitespace-pre-wrap leading-relaxed">
                                            {researchMutation.data.summary}
                                        </div>
                                    )}
                                    {researchMutation.error && (
                                        <div className="p-4 bg-red-950/20 border border-red-800 text-red-300 rounded-lg text-xs font-mono">
                                            Error: {researchMutation.error.message}
                                        </div>
                                    )}
                                </div>
                            </CardContent>
                        </Card>

                        {/* Coder Section */}
                        <Card className="bg-zinc-900 border-zinc-800 p-6 flex flex-col gap-4">
                            <CardHeader className="p-0">
                                <CardTitle className="text-lg font-bold text-purple-400 flex items-center gap-2">
                                    <Sparkles className="h-5 w-5" /> Coder Agent
                                </CardTitle>
                            </CardHeader>
                            <CardContent className="p-0 flex flex-col gap-4 flex-1">
                                <div className="flex gap-2">
                                    <Input
                                        type="text"
                                        className="flex-1 bg-black border-zinc-800 text-white outline-none placeholder:text-zinc-650"
                                        placeholder="Coding task (e.g. 'Write a test file for helper/date.ts')..."
                                        value={coderTask}
                                        onChange={(e) => setCoderTask(e.target.value)}
                                        onKeyDown={(e) => e.key === 'Enter' && handleCode()}
                                    />
                                    <Button
                                        onClick={handleCode}
                                        disabled={coderMutation.isPending}
                                        className="bg-purple-600 hover:bg-purple-500 text-white font-semibold"
                                    >
                                        Code
                                    </Button>
                                </div>

                                <div className="flex-1 overflow-y-auto">
                                    {coderMutation.isPending && (
                                        <div className="flex justify-center p-8"><Loader2 className="animate-spin text-zinc-500" /></div>
                                    )}
                                    {coderMutation.data && (
                                        <div className="bg-black p-4 border border-zinc-800 rounded-lg text-sm text-zinc-300">
                                            <h4 className="text-green-400 font-bold mb-2">Task Complete!</h4>
                                            <p className="font-mono text-xs text-zinc-400">Files Changed: {coderMutation.data.filesChanged?.join(', ') || 'None'}</p>
                                            <p className="mt-2 leading-relaxed text-zinc-300">{coderMutation.data.reasoning}</p>
                                        </div>
                                    )}
                                    {coderMutation.error && (
                                        <div className="p-4 bg-red-950/20 border border-red-800 text-red-300 rounded-lg text-xs font-mono">
                                            Error: {coderMutation.error.message}
                                        </div>
                                    )}
                                </div>
                            </CardContent>
                        </Card>
                    </div>
                </TabsContent>
            </Tabs>
        </div>
    );
}
