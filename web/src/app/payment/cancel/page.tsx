"use client";

import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";
import { Suspense } from "react";

function PaymentCancelContent() {
  const router = useRouter();

  return (
    <main className="min-h-screen bg-gradient-to-b from-white to-gray-50 flex flex-col items-center justify-center p-4">
      <div className="bg-white p-8 rounded-2xl shadow-lg text-center max-w-md w-full">
        <div className="mb-6">
          <div className="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg
              className="w-8 h-8 text-red-500"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </div>
          <h1 className="text-2xl font-bold text-gray-900">
            Payment Cancelled
          </h1>
          <p className="text-gray-600 mt-2">
            The payment process was cancelled or failed.
          </p>
        </div>
        <Button
          className="w-full text-lg py-6"
          variant="outline"
          onClick={() => router.push("/")}
        >
          Return Home
        </Button>
      </div>
    </main>
  );
}

export default function PaymentCancelPage() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <PaymentCancelContent />
    </Suspense>
  );
}
