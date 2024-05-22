require "time"

class DayOneEntry
  attr_reader :tags, :text, :timestamp, :uuid, :photos

  def initialize(timestamp:, uuid:, text: "", tags: [], photos: [])
    @tags = ["DayOneJournal"] + tags
    @timestamp = timestamp
    @uuid = uuid
    @photos = photos
    @text = clean_text(text, photos)
  end

  def title
    text.split("\n").first.gsub(/^# /, "")
  end

  def body
    text.split("\n")[1..-1].join("\n").gsub(/^\n/, "")
  end

  def dayone_url
    "[DayOne](dayone://view?entryId=#{uuid})"
  end

  def to_logseq_format
    tags = self.tags.reject { |tag| body.include?("##{tag}") }
      .map { |tag| "##{tag}" }
      .join(" ")
    body_content = "\t- "
    body_content += body.gsub("\n", "\n\t- ")

    output = "- ## #{title} #{tags}\n"
    output += "#{body_content}\n" unless body.empty?
    output += "\t- #{dayone_url}\n"
    output += "\t- dayone-id:: #{uuid}\n"
    output
  end

  def clean_text(text, photos)
    text.gsub!(/!\[\]\(dayone-moment:\/\/([A-F0-9]+)\)/) do
      identifier = $1
      photo = photos.find { |photo| photo.identifier == identifier }
      photo ? photo.to_markdown : ""
    end
    text.gsub!(/!\[\]\(dayone-moment:\/\w*\/\w+\)/, "") # Remove other assets like audio or video
    text.delete!("\\")
    text.gsub!(/\n{2}/, "\n")
    text
  end

  def self.from_json(json, use_photos: false)
    if json["text"].nil?
      return
    end

    photos = if use_photos
      json["photos"]&.map { |photo| DayOnePhoto.from_json(photo) } || []
    else
      []
    end

    DayOneEntry.new(
      text: json["text"],
      tags: json["tags"] || [],
      uuid: json["uuid"],
      timestamp: Time.parse(json["creationDate"]).localtime("+02:00"),
      photos: photos
    )
  end
end
