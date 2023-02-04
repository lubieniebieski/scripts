# frozen_string_literal: true

# Container for a link along with its reference number
class Link
  attr_reader :name, :url, :reference_no

  def initialize(name, url, reference_no)
    @name = name
    @url = url
    @reference_no = reference_no
  end
end

# Converts inline links to reference links in markdown files
class MarkdownConverter
  def initialize(file_name)
    @file_name = file_name
    @links = []
  end

  def extract_links!
    reference_no = 0
    File.foreach(@file_name) do |line|
      @links += line.scan(/\[(.*?)\]\((.*?)\)/).map do |name, url|
        reference_no += 1
        Link.new(name, url, reference_no)
      end
    end
    @links.uniq!(&:url)
  end

  def update_file
    extract_links!

    content = File.read(@file_name)
    @links.each do |link|
      content.gsub!(/\[#{link.name}\]\(.*?\)/, "[#{link.name}][#{link.reference_no}]")
    end
    return if references.empty?

    content += "\n"
    content += references.join("\n")

    File.open(@file_name, 'w') { |f| f.write(content) }
  end

  def references
    @references ||= @links.map { |link| "[#{link.reference_no}]: #{link.url}" }
  end

  def self.process_files(files)
    files.each do |file_name|
      if File.directory?(file_name)
        process_directory(file_name)
      elsif !File.file?(file_name)
        puts "Error: file not found: #{file_name}"
        exit 1
      else
        new(file_name).update_file
      end
    end
  end

  def self.process_directory(directory)
    Dir.glob("#{directory}/**/*.md") do |file_name|
      new(file_name).update_file
    end
  end
end

file_names = ARGV

MarkdownConverter.process_files(file_names)
