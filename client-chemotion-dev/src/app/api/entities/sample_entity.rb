module Entities
  class SampleEntity < Entities::SampleAttrEntity
    expose :molecule, using: Entities::MoleculeEntity
    expose :container, using: Entities::ContainerEntity
    expose :tag
    expose :segments, using: Entities::SegmentEntity
    expose :residues
    expose :elemental_compositions, using: Entities::ElementalCompositionEntity

    expose :code_log, using: Entities::CodeLogEntity

    class Level0 < SampleEntity
      include SamplePolicySerializable
      include SampleLevelEntity
      define_restricted_methods_for_level(0)

      def molecule
        {
          molecular_weight: object.molecule.try(:molecular_weight),
          exact_molecular_weight: object.molecule.try(:exact_molecular_weight),
        }
      end
    end

    class Level1 < SampleEntity
      include SamplePolicySerializable
      include SampleLevelEntity
      define_restricted_methods_for_level(1)
    end

    class Level2 < SampleEntity
      include SamplePolicySerializable
      include SampleLevelEntity
      define_restricted_methods_for_level(2)

      def analyses
        object.analyses && object.analyses.map {|x| x['datasets'] = {:datasets => []} if x['datasets'].present? }
      end
    end

    class Level3 < SampleEntity
      include SamplePolicySerializable
      include SampleLevelEntity
      define_restricted_methods_for_level(3)
    end

    class Level10 < SampleEntity
      include SamplePolicySerializable
    end
  end
end
