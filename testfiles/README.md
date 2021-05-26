# Chemotion-Spectra extension test files

To test the functionality of the Chemotion spectra extension, the files in this folder can be used.

## Steps to test:

 - Create a new collection
 - Add a new sample using the canonical smiles identifier notation: `O=Cc1nc2ccccc2nc1C`
 - Add four new analyses and attach the files according to this list:
    - `EI Mass.RAW`: 1H NMR
    - `IR.dx`: infrared absorption spectroscopy (IR)
    - `VLL-084_11.dx`: 13C NMR
    - `VLL-084_10.zip`: 13C NMR
  - Save the sample and wait.

If everything works, preview images should be seen and a blue button "Spectra Editor" should appear on the analysis.
