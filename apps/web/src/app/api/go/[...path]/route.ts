import { type NextRequest, NextResponse } from "next/server";

const TN_KERNEL_BASE =
	process.env.TORMENTNEXUS_TN_KERNEL_URL || "http://127.0.0.1:7778";

/**
 * TN Kernel reverse proxy.
 *
 * Browser requests to /api/go/<path> are forwarded to the TN Kernel
 * at TN_KERNEL_BASE/<path>.  The path is passed through verbatim so
 * that callers can target any kernel endpoint — e.g. /api/go/health
 * → http://127.0.0.1:7778/health, or /api/go/api/mcp/status →
 * http://127.0.0.1:7778/api/mcp/status.
 */

function remapPath(path: string): string {
	const map: Record<string, string> = {
		"api/imports": "api/import/summary",
		"api/healer": "api/healer/history",
		"api/deerflow": "api/deerflow/status",
		"api/cold-archive": "api/memory/cold-archive",
		"api/cli-harnesses": "api/tools/detect-cli-harnesses",
		"api/cloud-dev": "api/clouddev/sessions",
		"api/browser": "api/browser/status",
		"api/browser-extension": "api/browser-extension/stats",
		"api/logs-metrics": "api/metrics/stats",
		"api/observability": "api/pulse/status",
		"api/mesh": "api/mesh/status",
		"api/runtime": "api/runtime/status"
	};
	if (map[path]) {
		return map[path];
	}
	if (
		!path.startsWith("api/") &&
		!path.startsWith("trpc/") &&
		!path.startsWith("health") &&
		!path.startsWith("version") &&
		!path.startsWith("well-known")
	) {
		return "api/" + path;
	}
	return path;
}

export async function GET(
	request: NextRequest,
	{ params }: { params: Promise<{ path: string[] }> },
) {
	const resolvedParams = await params;
	const pathSegments = resolvedParams.path.join("/");
	const targetURL = `${TN_KERNEL_BASE}/${remapPath(pathSegments)}`;

	try {
		const response = await fetch(targetURL, {
			headers: {
				accept: "application/json",
				...Object.fromEntries(
					Object.entries(request.headers).filter(([key]) =>
						["authorization", "cookie"].includes(key.toLowerCase()),
					),
				),
			},
			signal: AbortSignal.timeout(5000),
		});

		const contentType = response.headers.get("content-type") || "";
		if (contentType.includes("application/json")) {
			const data = await response.json();
			return NextResponse.json(data, { status: response.status });
		}
		// Non-JSON response — pass through as text
		const text = await response.text();
		return new NextResponse(text, {
			status: response.status,
			headers: { "content-type": contentType },
		});
	} catch (error) {
		const message =
			error instanceof Error ? error.message : "TN Kernel unreachable";
		return NextResponse.json(
			{ success: false, error: message, kernelURL: targetURL },
			{ status: 502 },
		);
	}
}

export async function POST(
	request: NextRequest,
	{ params }: { params: Promise<{ path: string[] }> },
) {
	const resolvedParams = await params;
	const pathSegments = resolvedParams.path.join("/");
	const targetURL = `${TN_KERNEL_BASE}/${remapPath(pathSegments)}`;

	try {
		let body: string | null = null;
		const contentType = request.headers.get("content-type") || "";
		if (contentType.includes("application/json")) {
			body = await request.text();
		}

		const response = await fetch(targetURL, {
			method: "POST",
			headers: {
				"content-type": contentType || "application/json",
			},
			body,
			signal: AbortSignal.timeout(10000),
		});

		const respContentType = response.headers.get("content-type") || "";
		if (respContentType.includes("application/json")) {
			const data = await response.json();
			return NextResponse.json(data, { status: response.status });
		}
		const text = await response.text();
		return new NextResponse(text, {
			status: response.status,
			headers: { "content-type": respContentType },
		});
	} catch (error) {
		const message =
			error instanceof Error ? error.message : "TN Kernel unreachable";
		return NextResponse.json(
			{ success: false, error: message, kernelURL: targetURL },
			{ status: 502 },
		);
	}
}
