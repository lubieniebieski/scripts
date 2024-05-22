class DayOnePhoto
  attr_reader :identifier, :md5, :type
  ASSETS_PATH = "../assets/dayone/photos/"

  def initialize(identifier:, md5:, type:)
    @identifier = identifier
    @md5 = md5
    @type = type
  end

  def filename
    "#{md5}.#{type}"
  end

  def to_markdown
    path = "#{ASSETS_PATH}/#{filename}"
    "![](#{path})"
  end

  def self.from_json(json)
    DayOnePhoto.new(
      identifier: json["identifier"],
      md5: json["md5"],
      type: json["type"]
    )
  end
end
