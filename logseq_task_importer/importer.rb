class RSSTaskImporter
  require 'rss'
  require 'open-uri'
  require 'date'
  require 'pry'

  FEED_URL = "YOUR_FEED_URL"

  def store_last_date(date)
    open('tmp/.rss_last_item', 'w') do |f|
      f.puts(date)
    end
  end

  def retrieve_last_date
    Time.parse(File.read('tmp/.rss_last_item'))
  end

  def run
    URI.open(FEED_URL) do |rss|
      feed = RSS::Parser.parse(rss)
      last_item_date = retrieve_last_date
      puts 'test'
      puts store_last_date(feed.items.first.pubDate)
      puts "\n"
      feed.items.each do |item|
        break if item.pubDate <= last_item_date

        puts "- TODO #{item.title} #slack"
        puts '  collapsed:: true'
        date = item.pubDate.strftime('%Y-%m-%d %a')
        puts "  SCHEDULED: <#{date}>"
        puts "  - #{item.description}"
        puts "  - [Open](#{item.link})"
      end
    end
  end
end

importer = RSSTaskImporter.new
importer.run
