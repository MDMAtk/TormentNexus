'use client';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

export default function RedirectPage() {
    const router = useRouter();
    useEffect(() => {
        router.replace('/dashboard/mcp?tab=tool-sets');
    }, [router]);
    return null;
}
