"use client";

import { useEffect } from "react";
import * as Sentry from "@sentry/browser";

let initializedDsn: string | null = null;

type SentryInitProps = {
  dsn?: string;
  environment?: string;
};

export function SentryInit({ dsn, environment }: SentryInitProps) {
  useEffect(() => {
    if (!dsn || initializedDsn === dsn) {
      return;
    }

    Sentry.init({
      dsn,
      environment,
      attachStacktrace: true,
      sendDefaultPii: false,
    });
    initializedDsn = dsn;
  }, [dsn, environment]);

  return null;
}
