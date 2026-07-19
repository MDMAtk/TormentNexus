"use client";

import Link from "next/link";
import {
	LayoutDashboard,
	Terminal,
	Book,
	Brain,
	Settings,
	Activity,
	Network,
	Users,
	Clock,
	Shield,
	Home,
	MessageSquare,
	Zap,
	Workflow,
	Cpu,
	Box,
	History as HistoryIcon,
} from "lucide-react";
import { cn } from "../lib/utils";

/* Sidebar kept for legacy route compatibility. 
   Primary navigation is now via the UnifiedDashboard tabs. */
export function Sidebar() {
	return (
		<div className="flex h-full w-64 flex-col bg-gray-900 border-r border-gray-800">
			<div className="flex h-16 items-center px-6">
				<h1 className="text-xl font-bold text-white">TORMENTNEXUS</h1>
			</div>
			<nav className="flex-1 space-y-1 px-3 py-4">
				<Link
					href="/"
					className="group flex items-center rounded-md px-3 py-2 text-sm font-medium bg-gray-800 text-white"
				>
					<Home className="mr-3 h-5 w-5" />
					Dashboard
				</Link>
			</nav>
		</div>
	);
}
