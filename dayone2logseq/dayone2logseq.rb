require_relative "lib/day_one_entry"
require_relative "lib/day_one_photo"
require_relative "lib/day_one_journal"
require "pry"

# Usage
filename = ARGV[0]
if filename.nil?
  puts "You need to pass a filename in argument. Ex: ruby dayone2logseq.rb Journal.json"
  exit
end

journal = DayOneJournal.load_from_file(filename)

logseq_journal_path = ""
# journal.save_to_logseq(logseq_journal_path, 30)
journal.clean_journal_entries_from_logseq(logseq_journal_path)

Dir.foreach(logseq_journal_path) do |file|
  next unless File.extname(file) == ".md"
  file_path = File.join(logseq_journal_path, file)
  if File.file?(file_path) && (File.zero?(file_path) || File.read(file_path).chars.uniq.length < 3)
    puts "Deleting #{file_path}"
    File.delete(file_path)
  end
end
