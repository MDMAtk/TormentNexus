import { type NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
	try {
		const TN_KERNEL = process.env.TORMENTNEXUS_TN_KERNEL_URL || "http://127.0.0.1:7778";
		const res = await fetch(`${TN_KERNEL}/api/shutdown`, {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({}),
		});

		if (!res.ok) {
			const text = await res.text();
			return NextResponse.json(
				{ success: false, error: `Kernel error: ${text}` },
				{ status: res.status },
			);
		}

		const data = await res.json();
		return NextResponse.json(data);
	} catch (error) {
		const message = error instanceof Error ? error.message : "Unknown error";
		return NextResponse.json(
			{ success: false, error: message },
			{ status: 500 },
		);
	}
}
