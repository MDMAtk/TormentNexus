"use client";
import { Sheet, SheetContent, SheetTrigger } from "@tormentnexus/ui";
import { Button, StreamStatus } from "@tormentnexus/ui";
import { Menu } from "lucide-react";
import { useState } from "react";

interface NavigationProps {
    versionLabel?: string;
}

export function Navigation({ versionLabel = 'dev' }: NavigationProps) {
    const [open, setOpen] = useState(false);
    const isCloud = typeof window !== 'undefined' && (window.location.hostname.includes('hypernexus') || window.location.search.includes('brand=hypernexus'));

    return (
        <nav className="w-full bg-white dark:bg-zinc-900 border-b border-zinc-200 dark:border-zinc-800 px-6 py-4 flex items-center justify-between sticky top-0 z-50">
            <div className="flex items-center gap-6">
                <div className={`text-xl font-bold ${isCloud ? "bg-gradient-to-r from-blue-600 to-cyan-500" : "bg-gradient-to-r from-blue-500 to-purple-500"} bg-clip-text text-transparent`}>
                    {isCloud ? "HYPERNEXUS" : "TORMENTNEXUS"}
                </div>
            </div>

            {/* Mobile Navigation (Minimalist - Status only) */}
            <div className="md:hidden">
                <Sheet open={open} onOpenChange={setOpen}>
                    <SheetTrigger asChild>
                        <Button variant="ghost" size="icon">
                            <Menu className="h-6 w-6" />
                        </Button>
                    </SheetTrigger>
                    <SheetContent side="left" className="w-[300px] sm:w-[400px]">
                        <div className="flex flex-col gap-4 mt-8 h-full justify-between">
                            <div>
                                <h3 className="text-sm font-semibold uppercase text-zinc-500 mb-4">SYSTEM STATUS</h3>
                                <StreamStatus />
                            </div>
                            <div className="pt-8 border-t border-zinc-200 dark:border-zinc-800">
                                <span className="text-xs text-zinc-400">v{versionLabel}</span>
                            </div>
                        </div>
                    </SheetContent>
                </Sheet>
            </div>

            <div className="hidden md:flex items-center gap-4">
                <StreamStatus />
                <div className="text-xs text-zinc-400 border-l border-zinc-800 pl-4">
                    v{versionLabel}
                </div>
            </div>
        </nav>
    );
}
