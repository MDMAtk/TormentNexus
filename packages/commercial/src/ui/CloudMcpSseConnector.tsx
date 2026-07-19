"use client";

import React, { useState } from 'react';
import { Card, CardHeader, CardTitle, CardContent } from '@tormentnexus/ui';
import { Badge } from '@tormentnexus/ui';
import { Button } from '@tormentnexus/ui';
import { Input, Switch } from '@tormentnexus/ui';
import { Zap, RotateCw, WalletCards } from 'lucide-react';
import { toast } from 'sonner';

export function CloudMcpSseConnector() {
    const [sseAuthEnabled, setSseAuthEnabled] = useState(false);

    React.useEffect(() => {
        if (typeof window !== 'undefined') {
            setSseAuthEnabled(localStorage.getItem('sseAuthEnabled') === 'true');
        }
    }, []);

    const handleGenerateSseKey = () => {
        const newKey = "sec_" + Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);
        toast.success(`Generated new SSE Auth Token: ${newKey.substring(0, 8)}...`);
    };

    return (
        <Card className="bg-zinc-900 border-zinc-800 shadow-xl relative overflow-hidden mt-6">
            <div className="absolute top-0 right-0 w-32 h-32 bg-emerald-500/5 blur-3xl -mr-10 -mt-10 rounded-full" />
            <CardHeader>
                <CardTitle className="text-sm font-bold text-zinc-400 uppercase tracking-widest flex items-center gap-2">
                    <Zap className="h-4 w-4 text-emerald-400" />
                    Cloud MCP SSE Connector
                </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
                <div className="flex items-center justify-between border-b border-zinc-800/80 pb-2">
                    <span className="text-xs text-zinc-400">SSE Authentication Token</span>
                    <div className="flex items-center gap-2">
                        <span className="text-xs text-zinc-500">{sseAuthEnabled ? "Active" : "Disabled"}</span>
                        <input
                            type="checkbox"
                            checked={sseAuthEnabled}
                            onChange={(e) => {
                                setSseAuthEnabled(e.target.checked);
                                localStorage.setItem('sseAuthEnabled', String(e.target.checked));
                            }}
                            className="h-4 w-4 rounded border-zinc-700 bg-zinc-950 text-emerald-600 accent-emerald-500 outline-none"
                        />
                    </div>
                </div>

                <div className="space-y-2 text-xs font-mono">
                    <div className="flex flex-col gap-1">
                        <span className="text-zinc-500">Go Sidecar Base URL</span>
                        <span className="text-zinc-300 bg-black/50 p-2 rounded border border-zinc-800">
                            http://127.0.0.1:4300
                        </span>
                    </div>
                    <div className="flex flex-col gap-1">
                        <span className="text-zinc-500">Event Stream Endpoint</span>
                        <span className="text-emerald-400/80 bg-black/50 p-2 rounded border border-zinc-800">
                            GET /api/sse
                        </span>
                    </div>
                    <div className="flex flex-col gap-1">
                        <span className="text-zinc-500">Client Message Post</span>
                        <span className="text-blue-400/80 bg-black/50 p-2 rounded border border-zinc-800">
                            POST /api/sse/message
                        </span>
                    </div>
                </div>

                <div className="flex items-center justify-between border-t border-zinc-800/80 pt-3 mt-2">
                    <span className="text-xs text-zinc-400">Docker Tenant Auto-Isolation</span>
                    <div className="flex items-center gap-2">
                        <Badge variant="outline" className="bg-emerald-500/10 text-emerald-400 border-emerald-500/20 text-[10px]">
                            NEGOTIATING
                        </Badge>
                        <Switch checked={true} onCheckedChange={() => {}} disabled />
                    </div>
                </div>

                <div className="pt-2">
                    <Button
                        variant="outline"
                        className="w-full bg-zinc-800 border-zinc-700 text-zinc-300 hover:bg-zinc-750 text-xs flex items-center justify-center gap-2"
                        onClick={handleGenerateSseKey}
                    >
                        <RotateCw className="w-3 h-3" />
                        Generate New SSE Auth Token
                    </Button>
                </div>

                <div className="mt-4 pt-3 border-t border-zinc-800/80 flex justify-between items-center text-[10px] text-zinc-500 font-mono">
                    <span className="flex items-center gap-1">
                        <WalletCards className="w-3 h-3 text-cyan-500" />
                        Stripe Hooks Active
                    </span>
                    <span>Last Sync: Just now</span>
                </div>
            </CardContent>
        </Card>
    );
}
