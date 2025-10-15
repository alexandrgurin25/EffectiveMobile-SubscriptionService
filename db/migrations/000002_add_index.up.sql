CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_service_name ON subscriptions(service_name);

CREATE INDEX idx_subscriptions_user_service ON subscriptions(user_id, service_name);
CREATE INDEX idx_subscriptions_summary_query ON subscriptions(user_id, service_name, start_date, end_date);
