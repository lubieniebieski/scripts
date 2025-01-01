require 'fileutils'
require 'optparse'
require 'date'
require 'yaml'

# Parse command-line options
options = {}
OptionParser.new do |opts|
  opts.banner = "Usage: obsidian-journal-cleaner.rb [options]"

  opts.on("-d", "--dry-run", "Perform a dry run without moving files") do |d|
    options[:dry_run] = d
  end

  opts.on("-r", "--remove-original", "Remove original files after moving") do |r|
    options[:remove_original] = r
  end
end.parse!

# Define source and destination directories
source_dir = 'EG.ObsidianVault/Journals'
destination_dir = 'EG.ObsidianVault/Journals/01 Daily'

# Iterate through all the files in the source directory
Dir.glob("#{source_dir}/[0-9][0-9][0-9][0-9]/**/*.md").each do |file|
  # Extract the date from the file path
  date_match = file.match(%r{(\d{4})/(\d{4}-\d{2}-\d{2})\.md})
  next unless date_match

  year = date_match[1]
  date_str = date_match[2]

  # Parse the date
  date = Date.parse(date_str)
  created_at = date.strftime('%Y-%m-%d %H:%M')
  weekly = date.strftime('%Y-W%U')
  monthly = date.strftime('%Y-%m-M')
  yearly = date.strftime('%Y-Y')
  date_short = date.strftime('%A %d. %B')
  aliases = "#{date_short} #{year}"
  journal_start_date = date_str
  journal_end_date = date_str
  month = date.strftime('%m')

  # Define the properties to add
  properties = {
    'date' => date_str,
    'createdAt' => created_at,
    'weekly' => weekly,
    'monthly' => [monthly],
    'yearly' => yearly,
    'dateShort' => date_short,
    'aliases' => [aliases],
    'journal' => 'personal',
    'journal-start-date' => journal_start_date,
    'journal-end-date' => journal_end_date,
    'journal-section' => 'day'
  }

  # Convert properties to YAML front matter
  yaml_front_matter = properties.to_yaml
  yaml_front_matter = "#{yaml_front_matter}---\n"

  # Construct the new path with month subdirectory
  new_path = "#{destination_dir}/#{year}/#{month}/#{date_str}.md"

  if options[:dry_run]
    # Print the source and destination paths
    puts "Would move: #{file} -> #{new_path}"
  else
    # Create the destination directory if it doesn't exist
    FileUtils.mkdir_p(File.dirname(new_path))

    # Read the content of the file
    content = File.read(file)

    # Prepend the YAML front matter
    updated_content = "#{yaml_front_matter}\n## Summary\n#{content}"

    # Write the updated content to the new path
    File.write(new_path, updated_content)

    # Remove the original file if the option is set
    if options[:remove_original]
      FileUtils.rm(file)
    end

    puts "Moved and updated: #{file} -> #{new_path}"
  end
end