'use client';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

export default function RedirectPage() {
    const router = useRouter();
    useEffect(() => {
        router.replace('/dashboard?tab=cloud-orchestrator');
    }, [router]);
    return null;
}
