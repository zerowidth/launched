task :migrate do
  require "redis"
  require "json"

  if File.exist?("exported.json")
    data = JSON.parse(File.read("exported.json"))
  else
    data = {}
    from_redis = Redis.new(url: ENV.fetch("FROM_REDIS"), ssl_params: { verify_mode: OpenSSL::SSL::VERIFY_NONE })
    from_redis.scan_each(match: "launchd_plist:*", count: 1000, type: :hash).each do |key|
      puts "loading #{key}"
      plist = from_redis.hgetall(key)
      data[key] = plist
    end
    File.write("exported.json", JSON.dump(data))
  end

  puts "loaded #{data.size} keys"

  to_redis = Redis.new(url: ENV.fetch("TO_REDIS"), ssl_params: { verify_mode: OpenSSL::SSL::VERIFY_NONE })

  data.each_slice(50) do |slice|
    to_redis.pipelined do
      slice.each do |key, values|
        puts "writing #{key}"
        to_redis.hset(key, values)
      end
    end
  end
end
