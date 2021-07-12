require "helper"

module Cdx
  class Creator
    include Cdx::Helper
    attr_reader :doc, :str, :bond_length
    def initialize(args)
      @doc = args[:cdxml]
      @str = CdxStatic.init

      @tmp_id = 1000
    end

    def to_cdx
      lines = doc.split("\n")
      lines[3..-1].each do |line|
        translation(line)
      end
      raw_cdx = str_to_cdx(str.gsub("\n", " ").split(" "))
    end

    private

    def get_tmp_id
      @tmp_id += 1
    end

    def translation(input)
      line  = Nokogiri::XML(input)
      if line.root
        id = root_id(line.root)
        name = line.root.name
        if name == "CDXML"
          bond_length(line.root)
        elsif name == "page"
          @str += CdxStatic.page
        elsif name == "fragment"
          @str += CdxStatic.fragment(get_tmp_id)
        elsif name == "n"
          bid = backup_id
          @str += CdxNode.new(line.root, id, bid).content
        elsif name == "b"
          @str += CdxBond.new(line.root, id).content
        elsif name == "arrow"
          @str += CdxArrow.new(line.root, id).content
        elsif name == "t"
          @str += CdxStr.new(line.root, id).content
        end
      else
        @str += CdxStatic.ending
      end
    end

    def bond_length(root)
      @bond_length = root["BondLength"].to_i
    end

    def root_id(root)
      root["id"] = get_tmp_id if !root["id"]

      hex = "#{"%04X" % root["id"].to_i}"
      self.little_endian(hex) + "00 00 "
    end

    def backup_id
      empty_root = { "id": nil }
      root_id(empty_root)
    end

    def str_to_cdx(input)
      hex = input.each_with_object("") { |inp, out| out << str_to_hex(inp) }
      hex_to_cdx(hex)
    end

    def str_to_hex(input)
      begin
        Integer("0X#{input}")
      rescue
        Integer("0X#{0}")
      end
    end

    def hex_to_cdx(input)
      input.encode!(Encoding::ISO_8859_1)
    end

    def save(data, path)
      File.binwrite(path, data)
    end
  end
end
