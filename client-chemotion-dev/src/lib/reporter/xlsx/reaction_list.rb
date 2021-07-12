module Reporter
  module Xlsx
    class ReactionList
      IMG_HEIGHT = 51
      IMG_WIDTH = 180
      TEXT_SIZE = 14
      ROW_PRD_BEGIN = 3
      C_W = 25
      IMG_BEGIN_X = 1
      IMG_BEGIN_Y = 3
      ROW_HEIGHT = 40

      def initialize(args)
        @objs = args[:objs]
        @mol_serials = args[:mol_serials] || []
      end

      def create(file_name)
        Axlsx::Package.new do |p|
          p.workbook.add_worksheet(name: 'ELN Report Reaction List') do |st|
            @sheet = st

            row_main_title
            row_info
            row_sub_title
            row_content
            col_widths

            @sheet['A3:H#{row_counts}'].each do |c| c.style = Axlsx::STYLE_THIN_BORDER end
          end
          p.serialize(file_name)
        end
      end

      private

      def row_counts
        @objs.length + ROW_PRD_BEGIN
      end

      def col_widths
        @sheet.column_widths C_W, C_W, C_W, C_W, C_W, C_W, C_W, C_W
      end

      def row_main_title
        @sheet.add_row [
          'Article Reference', '', '', '', '', 'Article DOI', '', ''
        ], sz: TEXT_SIZE, height: ROW_HEIGHT, b: true
      end

      def row_info
        @sheet.add_row [
          'reference', '', '', '', '', 'doi:', '', ''
        ], sz: TEXT_SIZE, height: ROW_HEIGHT
        @sheet['A2:H2'].each do |c| c.color = 'aaaaaa' end
      end

      def row_sub_title
        @sheet.add_row [
          'Label', 'Image', 'Name', 'InChI', 'InChIKey', 'Long-RInChIKey',
          'Web-RInChIKey', 'Short-RInChIKey'
        ], sz: TEXT_SIZE, height: ROW_HEIGHT, b: true
      end

      def add_content_to_row(p, long_key, web_key, short_key)
        serial = Reporter::Helper.mol_serial(p[:molecule][:id], @mol_serials)
        @sheet.add_row [
          serial,
          '',
          p[:showed_name],
          p[:molecule][:inchistring],
          p[:molecule][:inchikey],
          long_key,
          web_key,
          short_key,
          ''
        ], sz: TEXT_SIZE, height: ROW_HEIGHT,
           types: %i[text string string string string string string string string]
      end

      def add_img_to_row(p)
        img_src = Reporter::Helper.mol_img_path(p)
        return if img_src.nil?
        @sheet.add_image(image_src: img_src) do |img|
          img.height = IMG_HEIGHT.to_i
          img.width = IMG_WIDTH.to_i
          img.start_at IMG_BEGIN_X, (IMG_BEGIN_Y + @counter)
        end
      end

      def row_content
        @counter = 0
        @objs.each do |obj|
          long_key, web_key, short_key = Reporter::Helper.get_rinchi_keys(obj)

          obj[:products].each do |p|
            add_content_to_row(p, long_key, web_key, short_key)
            add_img_to_row(p)
            @counter += 1
          end
        end
      end
    end
  end
end
