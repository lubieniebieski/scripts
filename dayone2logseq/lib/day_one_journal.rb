require "json"

class DayOneJournal
  attr_reader :entries

  def initialize(filename, use_photos: false)
    @entries = JSON.parse(File.read(filename))["entries"].map { |entry| DayOneEntry.from_json(entry, use_photos: use_photos) }.compact
  end

  def by_date(count = 30)
    entries.last(count).group_by { |entry| entry.timestamp.strftime("%Y_%m_%d") }
  end

  def self.load_from_file(filename, use_photos: false)
    DayOneJournal.new(filename, use_photos: use_photos)
  end

  def find(uuid)
    entries.find { |entry| entry.uuid == uuid }
  end

  def clean_journal_entries_from_logseq(logseq_journal_path)
    Dir.glob("#{logseq_journal_path}/*.md").each do |file|
      content = File.read(file)
      entries_uuids = content.scan(/dayone-id:: ([A-F0-9-]+)/).flatten
      entries_uuids.each do |uuid|
        if find(uuid)
          content.gsub!(/- ## [^\n]*\n(\t- [^\n]*\n)*\t- dayone-id:: #{uuid}\n/, "")
          File.write(file, content)
        end
      end
    end
  end

  def save_to_logseq(logseq_journal_path, entries_count = 5)
    by_date(entries_count).each do |date, entries|
      File.open("#{logseq_journal_path}/#{date}.md", "a+") do |file|
        existing_content = file.read

        entries.each do |entry|
          if existing_content.include?(entry.uuid)
            next
          elsif !existing_content.empty?
            file.puts
          end
          file.puts entry.to_logseq_format
        end
      end
    end
  end
end
