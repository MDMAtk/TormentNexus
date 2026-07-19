"use client";

import { CommercialView } from "@/components/commercial/CommercialView";

export default function CommercialPage() {
	return (
		<div className="max-w-7xl mx-auto">
			<div className="flex justify-between items-center mb-8">
				<div>
					<h1 className="text-3xl font-bold tracking-tight">
						Commercial Management
					</h1>
					<p className="text-muted-foreground mt-2">
						Governance, compliance, and distributed orchestration for
						TORMENTNEXUS.
					</p>
				</div>
			</div>
			<CommercialView />
		</div>
	);
}
