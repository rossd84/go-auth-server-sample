-- Enable pg_cron for automated removal of refresh_tokens every day at midnight.

CREATE EXTENSION IF NOT EXISTS pg_cron;

-- Schedule a daily cleanup at 2:00 AM
SELECT cron.schedule(
  'delete_expired_tokens',
  '0 2 * * *',
  $$DELETE FROM refresh_tokens WHERE expires_at < NOW()$$
);
