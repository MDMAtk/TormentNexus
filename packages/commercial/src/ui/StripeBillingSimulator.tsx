"use client";

import React, { useState } from 'react';
import { Card, CardHeader, CardTitle, CardContent } from '@tormentnexus/ui';
import { Badge } from '@tormentnexus/ui';
import { Button } from '@tormentnexus/ui';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@tormentnexus/ui';
import { WalletCards, ExternalLink } from 'lucide-react';

export function StripeBillingSimulator() {
    const [checkoutOpen, setCheckoutOpen] = useState(false);
    const [billingPortalOpen, setBillingPortalOpen] = useState(false);

    return (
        <>
            <Card className="bg-zinc-900 border-zinc-800 shadow-xl relative overflow-hidden">
                <div className="absolute top-0 right-0 w-32 h-32 bg-cyan-500/5 blur-3xl -mr-10 -mt-10 rounded-full" />
                <CardHeader>
                    <CardTitle className="text-sm font-bold text-zinc-400 uppercase tracking-widest flex items-center gap-2">
                        <WalletCards className="h-4 w-4 text-cyan-400" />
                        HyperNexus Cloud Billing (Stripe)
                    </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                    <div className="flex items-center justify-between border-b border-zinc-800/80 pb-3">
                        <div>
                            <div className="text-xs text-zinc-500 uppercase">Subscription Plan</div>
                            <div className="text-base font-bold text-white mt-1">Commercial Cloud SaaS</div>
                        </div>
                        <Badge variant="outline" className="bg-cyan-500/10 text-cyan-400 border-cyan-500/20 text-xs">
                            ACTIVE (PAID)
                        </Badge>
                    </div>
                    <div className="grid grid-cols-2 gap-4 text-xs font-mono">
                        <div>
                            <div className="text-zinc-500">Monthly Price</div>
                            <div className="text-zinc-200 mt-1">$499.00 / month</div>
                        </div>
                        <div>
                            <div className="text-zinc-500">Next Invoice Date</div>
                            <div className="text-zinc-200 mt-1">July 25, 2026</div>
                        </div>
                        <div>
                            <div className="text-zinc-500">Payment Source</div>
                            <div className="text-zinc-200 mt-1">Visa ending in 4242</div>
                        </div>
                        <div>
                            <div className="text-zinc-500">Customer ID</div>
                            <div className="text-zinc-400 mt-1">cus_R8vB42tX910a</div>
                        </div>
                    </div>
                    <div className="flex gap-2 pt-2">
                        <Button
                            className="bg-cyan-600 hover:bg-cyan-500 text-white border-transparent text-xs"
                            onClick={() => setBillingPortalOpen(true)}
                        >
                            Manage via Stripe Portal
                        </Button>
                        <Button
                            variant="outline"
                            className="bg-zinc-800 border-zinc-700 text-zinc-300 hover:bg-zinc-750 text-xs"
                            onClick={() => setCheckoutOpen(true)}
                        >
                            Upgrade Plan
                        </Button>
                    </div>
                </CardContent>
            </Card>

            <Dialog open={checkoutOpen} onOpenChange={setCheckoutOpen}>
                <DialogContent className="sm:max-w-md bg-zinc-950 border-zinc-800 text-zinc-200">
                    <DialogHeader>
                        <DialogTitle className="flex items-center gap-2">
                            <WalletCards className="h-5 w-5 text-cyan-400" />
                            Stripe Checkout Simulator
                        </DialogTitle>
                        <DialogDescription className="text-zinc-400">
                            Configure or upgrade your HyperNexus Cloud Subscription.
                        </DialogDescription>
                    </DialogHeader>
                    <div className="py-4 space-y-4">
                        <div className="rounded-lg bg-zinc-900 border border-zinc-800 p-4 text-xs space-y-2">
                            <div className="font-semibold text-white">HyperNexus Pro Plan Upgrade</div>
                            <div className="text-zinc-400">Unlimited scale, SOC 2 compliance logging, and dedicated SLA support.</div>
                            <div className="text-sm font-bold text-cyan-400 pt-1">$999.00 / month</div>
                        </div>
                        <div className="text-xs text-zinc-500 italic text-center">
                            Simulating secure transaction redirect to stripe.com...
                        </div>
                    </div>
                    <DialogFooter>
                        <Button
                            variant="outline"
                            onClick={() => setCheckoutOpen(false)}
                            className="bg-zinc-900 border-zinc-800 hover:bg-zinc-800 text-zinc-300"
                        >
                            Cancel
                        </Button>
                        <Button
                            onClick={() => {
                                setCheckoutOpen(false);
                            }}
                            className="bg-cyan-600 hover:bg-cyan-500 text-white border-transparent"
                        >
                            Complete Payment ($999.00/mo)
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

            <Dialog open={billingPortalOpen} onOpenChange={setBillingPortalOpen}>
                <DialogContent className="sm:max-w-lg bg-zinc-950 border-zinc-800 text-zinc-200">
                    <DialogHeader>
                        <DialogTitle className="flex items-center gap-2">
                            <ExternalLink className="h-5 w-5 text-cyan-400" />
                            Stripe Customer Portal (Simulated)
                        </DialogTitle>
                        <DialogDescription className="text-zinc-400">
                            Update payment methods, view invoices, or cancel subscriptions securely.
                        </DialogDescription>
                    </DialogHeader>
                    <div className="py-4 space-y-4 text-xs">
                        <div className="rounded-lg border border-zinc-800 bg-zinc-900 p-4 space-y-3">
                            <div className="flex justify-between items-center">
                                <span className="font-semibold text-white">Payment Method</span>
                                <span className="text-zinc-400">Visa **** 4242 (Expires 12/29)</span>
                            </div>
                            <div className="flex justify-between items-center">
                                <span className="font-semibold text-white">Billing Address</span>
                                <span className="text-zinc-400">100 Pine St, San Francisco, CA</span>
                            </div>
                        </div>
                        <div className="space-y-2">
                            <div className="font-semibold text-zinc-400 uppercase tracking-wider text-[10px]">Invoice History</div>
                            <div className="divide-y divide-zinc-800 border border-zinc-800 rounded bg-zinc-900">
                                <div className="flex justify-between items-center p-3">
                                    <span>June 25, 2026</span>
                                    <span className="font-mono">$499.00 (Paid ✓)</span>
                                </div>
                                <div className="flex justify-between items-center p-3">
                                    <span>May 25, 2026</span>
                                    <span className="font-mono">$499.00 (Paid ✓)</span>
                                </div>
                            </div>
                        </div>
                    </div>
                    <DialogFooter>
                        <Button
                            onClick={() => setBillingPortalOpen(false)}
                            className="bg-cyan-600 hover:bg-cyan-500 text-white border-transparent w-full"
                        >
                            Return to HyperNexus
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </>
    );
}
