"use client";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex min-h-screen bg-black">
      <main className="flex-1 overflow-auto min-w-0">
        {children}
      </main>
    </div>
  );
}
