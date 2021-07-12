class CreateCollectorErrors < ActiveRecord::Migration
  def change
    create_table :collector_errors do |t|
      t.string :error_code

      t.timestamps null: false
    end
  end
end
