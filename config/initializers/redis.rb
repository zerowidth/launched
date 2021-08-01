REDIS = ConnectionPool.new(size: ENV.fetch("RAILS_MAX_THREADS", 5), timeout: 5) do
  Redis.new(url: ENV["REDIS_URL"], port: ENV["REDIS_PORT"], db: ENV["REDIS_DB"])
end
