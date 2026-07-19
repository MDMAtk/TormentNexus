"use client";

import React, { useState } from 'react';
import { Card, CardHeader, CardTitle, CardContent } from '@tormentnexus/ui';
import { Shield } from 'lucide-react';
import { toast } from 'sonner';

export function CorporateModelFallback() {
    const [corporateIsolation, setCorporateIsolation] = useState(false);

    React.useEffect(() => {
        if (typeof window !== 'undefined') {
            setCorporateIsolation(localStorage.getItem('corporateIsolation') === 'true');
        }
    }, []);

    const toggleCorporateIsolation = (checked: boolean) => {
        setCorporateIsolation(checked);
        localStorage.setItem('corporateIsolation', String(checked));
        if (checked) {
            toast.success("Corporate Local Model Isolation enabled. External providers are now restricted.");
        } else {
            toast.info("Corporate Local Model Isolation disabled. External providers restored.");
        }
    };

    return (
        <Card className="bg-zinc-900 border-zinc-800 shadow-xl relative overflow-hidden">
            <div className="absolute top-0 right-0 w-32 h-32 bg-amber-500/5 blur-3xl -mr-10 -mt-10 rounded-full" />
            <CardHeader>
                <CardTitle className="text-sm font-bold text-zinc-400 uppercase tracking-widest flex items-center gap-2">
                    <Shield className="h-4 w-4 text-amber-500" />
                    Corporate Model Fallback Configuration
                </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
                <div className="flex items-center justify-between border-b border-zinc-800/80 pb-3">
                    <div>
                        <div className="text-xs text-zinc-500 uppercase">Isolation Strategy</div>
                        <div className="text-sm font-bold text-white mt-1">Local Compliance Restriction</div>
                    </div>
                    <div className="flex items-center gap-2">
                        <span className="text-xs text-zinc-400">Isolation Active</span>
                        <input
                            type="checkbox"
                            checked={corporateIsolation}
                            onChange={(e) => toggleCorporateIsolation(e.target.checked)}
                            className="h-4 w-4 rounded border-zinc-700 bg-zinc-950 text-amber-600 accent-amber-500 outline-none"
                        />
                    </div>
                </div>

                <div className="space-y-3">
                    <div className="bg-amber-500/10 border border-amber-500/20 p-3 rounded text-xs text-amber-200/90 leading-relaxed">
                        When Isolation Strategy is active, all public API calls are blocked.
                        The fallback chain will strictly resolve to local, air-gapped endpoints.
                    </div>
                    <div className="text-xs font-mono bg-black/50 p-3 rounded border border-zinc-800 text-zinc-400">
                        <span className="text-zinc-500">1.</span> gemma-4-e2b (Ollama)<br/>
                        <span className="text-zinc-500">2.</span> llama-3-8b-instruct (vLLM)<br/>
                        <span className="text-zinc-500">3.</span> mistral-7b-local (LM Studio)
                    </div>
                </div>
            </CardContent>
        </Card>
    );
}
