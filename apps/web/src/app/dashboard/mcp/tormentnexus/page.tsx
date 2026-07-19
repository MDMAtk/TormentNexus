'use client';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

export default function RedirectPage() {
    const router = useRouter();
    useEffect(() => {
        router.replace('/dashboard/mcp?tab=tormentnexus');
    }, [router]);
    return null;
}
