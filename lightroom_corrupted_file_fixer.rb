#!/usr/bin/env ruby
require "fileutils"

BACKUP_FOLDER = "/Volumes/Photos/USB import"
corrupted_files_paths = IO.read('./corrupted photos.txt').split("\n")[1..-1]

puts "Replacing corrupted files from backup"
corrupted_files_paths.each do |path|
  path = path.strip
  puts "### " + File.basename(path)
  backup_copy = Dir.glob("#{BACKUP_FOLDER}/**/#{File.basename(path)}").first
  if backup_copy.nil?
    puts "\t NOT FOUND"
    next
  end
  FileUtils.cp(backup_copy, File.dirname(path))
  puts "\tOK"
end
