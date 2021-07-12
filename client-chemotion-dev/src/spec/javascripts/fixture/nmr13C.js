const nmrData13C = {
  '2462': {
    formula: 'C16H14N2O2',
    content: {
      "ops": [
        {
          "insert": "13C NMR (101 MHz, DMSO, ppm) δ = 154.0 (Cq, 1C), 153.1 (Cq, 1C), 135.9 (Cq, 1C), 133.2 (Cq, 1C), 132.5 (Cq, 1C), 130.4 (+, CH, 1C), 130.1 (+, CH, 1C), 129.7 (+, CH, 1C), 129.3 (+, CH, 2C), 127.8 (+, CH, 2C), 123.4 (+, CH, 1C), 115.2 (+, CH, 1C), 57.9 (+, CH3, 1C), 44.5 (+, CH3, 1C).\n"
        }
      ]
    },
    expected: '',
  },
  'REPO-42343-Nicolai': {
    formula: 'C12H21N5O2',
    content: {"ops":[{"attributes":{"script":"super"},"insert":"13"},{"insert":"C NMR (100 MHz, CDCl"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":", ppm) δ = 163.5, 152.3, 142.7 (+, CH), 104.1, 59.9 (-, CH"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 50.7 (+, CH), 48.1 (+, CH), 23.7 (+, 2C, CH"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":"), 19.1 (+, 2C, CH"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":"), 14.6 (+, CH"},{"attributes":{'script':"sub"},"insert":"3"},{"insert":").\n"}]},
    expected: '',
  },
  'REPO-41965-Claudine': {
    formula: 'C30H46N4O5',
    content: {"ops":[{"attributes":{"script":"super"},"insert":"13"},{"insert":"C NMR (100 MHz, chloroform-d"},{"attributes":{"script":"sub"},"insert":"1"},{"insert":") δ = 175.1 (C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":", "},{"attributes":{"italic":true},"insert":"C"},{"insert":"ONH), 174.4 (C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":", "},{"attributes":{"italic":true},"insert":"C"},{"insert":"ONH), 173.3 (C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":", "},{"attributes":{"italic":true},"insert":"C"},{"insert":"ONH), 171.4 (C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":", "},{"attributes":{"italic":true},"insert":"C"},{"insert":"ONH), 154.0 (C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":", tyr-CH"},{"attributes":{"italic":true},"insert":"C"},{"insert":"O), 131.7 (C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":", tyr-"},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"ar"},{"insert":"), 129.2 (+, 2 × tyr-"},{"attributes":{"italic":true},"insert":"C"},{"insert":"H), 124.1 (+, 2 × tyr-"},{"attributes":{"italic":true},"insert":"C"},{"insert":"H), 78.3 (Cq, tyr-O"},{"attributes":{"italic":true},"insert":"C"},{"insert":"(CH"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":")"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":"), 62.9 (+, tyr-"},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"α"},{"insert":"H), 57.7 (+, pro-"},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"α"},{"insert":"H), 57.7 (+, ile-"},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"α"},{"insert":"H), 54.1 (+, nle-"},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"α"},{"insert":"H), 46.8 (–, pro-"},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"δ"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 34.1 (–, tyr-"},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 33.6 (+, ile-"},{"attributes":{"italic":true},"insert":"C"},{"insert":"H), 28.7 (–, nle-"},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"β"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 28.6 (+, 2 × tyr-"},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":"), 27.5 (–, ile-"},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 24.8 (–, pro-"},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"β"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 24.7 (–, pro-"},{"attributes":{"italic":true},"insert":"C"},{"insert":"γH"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 25.6 (–, nle-"},{"attributes":{"italic":true},"insert":"C"},{"insert":"γH"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 22.2 (–, nle-"},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"δ"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 15.6 (+, ile-"},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":"), 13.8 (+, nle-"},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":"), 10.5 (+, ile-"},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":") ppm.\n"},{"attributes":{"align":"justify"},"insert":"\n"},{"insert":"\n"}]},
    expected: ' count: 29/30',
  },
  'REPO-42299-Klein': {
    formula: 'C14H9CIN2',
    content: {"ops":[{"attributes":{"script":"super"},"insert":"13"},{"insert":"C NMR (100 MHz, CDCl"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":", ppm) δ = 128.3, 128.5 (2C), 129.4, 129.8 (2C), 129.9, 130.6, 131.0, 136.9, 141.2, 141.2, 146.3, 153.2.\n"}]},
    expected: '',
  },
  'REPO-42286 Robin': {
    formula: 'C18H18S2',
    content: {"ops":[{"attributes":{"script":"super"},"insert":"13"},{"insert":"C NMR (100 MHz, CDCl"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":", ppm) δ = 138.2 (C"},{"attributes":{"script":"sub"},"insert":"quat"},{"insert":", "},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"CH"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":"), 134.2 (+, 2x"},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 133.9 (+, 2x"},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 133.4 (C"},{"attributes":{"script":"sub"},"insert":"quat"},{"insert":", "},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"S), 129.8 (+, 2x"},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 129.6 (C"},{"attributes":{"script":"sub"},"insert":"quat"},{"insert":", "},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"S), 123.0 (+, 2x"},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 128.0 (+, "},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 57.4 (–, 3x"},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 42.8 (C"},{"attributes":{"script":"sub"},"insert":"quat"},{"insert":", "},{"attributes":{"italic":true},"insert":"C"},{"insert":"CH"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 42.6 (C"},{"attributes":{"script":"sub"},"insert":"quat"},{"insert":", "},{"attributes":{"italic":true},"insert":"C"},{"insert":"CH"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 21.3 (+, "},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":").\n"}]},
    expected: '',
  },
  'REPO-41987-Nicolai': {
    formula: 'C26H40O3',
    content: {"ops":[{"attributes":{"script":"super"},"insert":"13"},{"insert":"C NMR (126 MHz, DMSO-d"},{"attributes":{"script":"sub"},"insert":"6"},{"insert":") δ [ppm] = 174.9 (C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":", COOH), 81.2 (+, CH-6), 57.3 (+, CH), 56.0 (+, OCH"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":"), 55.5 (+, CH), 47.4 (+, CH), 42.9 (C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":"), 42.4 (C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":"), 39.2 (+, CH"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 38.6 (+, CH), 34.9 (C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":"), 34.9(-, CH"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 32.8(-, CH"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 30.2 (+, CH), 29.7 (+, CH), 27.2 (-, CH"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 24.6 (-, CH"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 23.8 (-, CH"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 22.3 (-, CH"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 21.9 (+, CH), 20.7 (+, CH), 19.7 (+, CH"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":"-21), 19.3 (+, CH"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":"-19), 12.8 (-, CH"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"-4), 12.1 (-, CH"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 12.0 (+, CH"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":"-18).\n"}]},
    expected: '',
  },
  'REPO-2536-Daniel': {
    formula: 'C21H19N',
    content: {"ops":[{"attributes":{"bold":true},"insert":"13C NMR "},{"insert":"(101 MHz, CDCl3) "},{"attributes":{"italic":true},"insert":"δ"},{"insert":" [ppm] = 144.1 (Cquat .), 141.1 (Cquat . ), 139.6 (Cquat. ), 137.4 (Cquat.), 137.1 (Cquat.), 134.7 (+, CArH), 134.7 (+, CArH), 133.7 (+, CArH), 133.6 (+, CArH), 132.8 (+, CArH), 132.6 (+, CArH), 131.2 (+, CArH), 36.51 (–, CH2), 36.41 (–, CH2), 36.33 (–, CH2), 36.16 (–, CH2).\n"}]},
    expected: ' count: 16/21',
  },
  'REPO-41901-Simone': {
    formula: 'C16H23N',
    content: {"ops":[{"attributes":{"script":"super"},"insert":"13"},{"insert":"C NMR (100 MHz, CDCl"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":", ppm) δ = 14.1, 22.6, 27.0, 29.2, 29.2, 30.2, 31.8, 46.4, 77.0, 100.8, 109.3, 119.1, 120.9, 121.2, 127.7, 128.5, 135.9. \n"}]},
    expected: ' count: 17/16',
  },
  'REPO-41853-Steven': {
    formula: 'C6H10O2S2',
    content: {"ops":[{"attributes":{"script":"super"},"insert":"13"},{"insert":"C NMR (100 MHz, CDCl"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":", ppm) δ = 172.9, 52.3, 41.7, 31.4, 31.2 (2C).\n"}]},
    expected: '',
  },
  'REPO-2875-Fabian': {
    formula: 'C72H32F30N8O2',
    content: {"ops":[{"attributes":{"bold":true,"script":"super"},"insert":"13"},{"attributes":{"bold":true},"insert":"C NMR"},{"insert":" (101 MHz, CDCl"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":") δ 160.3, 155.1, 143.1, 140.2, 138.8, 138.1, 137.6, 126.8, 126.3, 125.9, 125.0, 124.3, 124.2, 124.1, 121.6, 121.5, 121.0, 120.9, 120.4, 119.9, 110.0, 109.2, 108.1, 29.9."},{"attributes":{"align":"justify"},"insert":"\n"},{"insert":"\n\n"}]},
    expected: ' count: 24/72',
  },
  'REPO-41931-Mareen': {
    formula: 'C26H28N4',
    content: {"ops":[{"attributes":{"script":"super"},"insert":"13"},{"insert":"C NMR (101 MHz, CDCl"},{"attributes":{"script":"sub"},"insert":"3,"},{"insert":" ppm) "},{"attributes":{"italic":true},"insert":"δ"},{"insert":" = 143.5 ("},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 142.9 ("},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 140.4 ("},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 139.6 ("},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 139.2 ("},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 137.9 ("},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 137.1 ("},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 135.1 ("},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 133.8 ("},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 133.8 ("},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 133.5 ("},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 133.5 ("},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 132.7 ("},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 132.5 ("},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 132.0 ("},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 128.7 ("},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 126.1 ("},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 116.7 ("},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 57.0 ("},{"attributes":{"italic":true},"insert":"C"},{"insert":"), 35.6 ("},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 35.5 ("},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 35.4 ("},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 33.8 ("},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 29.9 (3 × "},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":").\n"}]},
    expected: '',
  },
  'REPO-2516-Simone': {
    formula: 'C11H16OS2',
    content: {"ops":[{"insert":"13","attributes":{"script":"sup"}},{"insert":"C NMR (125 MHz, CDCl"},{"insert":"3","attributes":{"script":"sub"}},{"insert":", ppm), δ = 12.9, 25.7, 27.3, 37.2, 37.3, 38.3, 42.1, 44.7, 127.6, 131.2, 218.3.\n"}]},
    expected: '',
  },
  'REPO-2523-Simone': {
    formula: 'C14H18O2S2',
    content: {"ops":[{"attributes":{"script":"super"},"insert":"13"},{"insert":"C NMR (100 MHz, CDCl"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":", ppm) δ = 2x26.4 (2C, 2 dia*), 37.1 (2C), 38.4 (2C), 40.8 (2C), 2x41.8 (2C, 2 dia*), 2x125.8 (1C, 2 dia*), 2x132.2 (1C, 2 dia*), 217.5+217.6 (2C, 2 dia*). *dia = 2 Cs as diastereomers. \n"}]},
    expected: '',
  },
  'REPO-42164-Hannes': {
    formula: 'C12H12N2',
    content: {"ops":[{"attributes":{"script":"super"},"insert":"13"},{"insert":"C NMR (126 MHz, CDCl"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":", ppm): δ = 158.8 (2C, 2- and 6-"},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":"), 150.5 (2C, 2'- and 6'-"},{"attributes":{"italic":true},"insert":"C"},{"insert":"H), 146.2 (4- or 4'-"},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":"), 146.1 (4- or 4'-"},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":"), 118.1 (4C, 3-, 5-, 3'- and 5'-"},{"attributes":{"italic":true},"insert":"C"},{"insert":"H) ,24.6 (2C, 2 x "},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":").\n"}]},
    expected: '',
  },
  'REPO-42129-Stefan': {
    formula: 'C19H24O3Si2',
    content: {"ops":[{"attributes":{"script":"super"},"insert":"13"},{"insert":"C NMR (100 MHz, CDCl"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":", ppm) δ = 190.2 (+, CHO), 165.0 (C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":", COO), 139.6 (C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":", C"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 135.0 (+, 2 x C"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"H), 133.4 (C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":", C"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 125.4 (C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":", 2 x C"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 103.6 (C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":", 2 x CCTMS), 100.5 (C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":", 2 x CCTMS), 52.9 (+, OCH"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":"), -0.2 (+, 6 x SiCH"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":").\n"}]},
    expected: '',
  },
  'REPO-42292-Robin': {
    formula: 'C19H20S2',
    content: {"ops":[{"attributes":{"script":"super"},"insert":"13"},{"insert":"C NMR (100 MHz, CDCl"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":", ppm) δ = 138.4 (C"},{"attributes":{"script":"sub"},"insert":"quat"},{"insert":", 2x"},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"CH"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 128.9 (+, 4x"},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 128.6 (+, 4x"},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 127.2 (+, 2x"},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 57.0 (–, 3xC"},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 41.3 (C"},{"attributes":{"script":"sub"},"insert":"quat"},{"insert":", 2x"},{"attributes":{"italic":true},"insert":"C"},{"insert":"CH"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":"), 36.2 (–, 2xPh"},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"2"},{"insert":").\n"}]},
    expected: '',
  },
  'REPO-2427-Sylvia': {
    formula: 'C11H16S2',
    content: {"ops":[{"insert":"13C NMR (100 MHz, CDCl3, ppm) δ = 17.5, 25.2, 28.9, 37.2, 37.6, 41.0, 45.2, 47.0, 76.7, 77.3, 123.3, 129.9, 210.8.\n"}]},
    expected: ' count: 13/11',
  },
  'REPO-42236-Robin': {
    formula: 'C12H9I',
    content: {"ops":[{"attributes":{"script":"super"},"insert":"13"},{"insert":"C NMR (100 MHz, CDCl"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":", ppm) δ = 146.8 (C"},{"attributes":{"script":"sub"},"insert":"quat"},{"insert":", "},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"C"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 144.4 (C"},{"attributes":{"script":"sub"},"insert":"quat"},{"insert":", C"},{"attributes":{"script":"sub"},"insert":"Ar"},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 139.6 (+, "},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 130.2 (+, "},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 129.4 (+, "},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 128.9 (+, "},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 128.2 (+, "},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 128.1 (+, "},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 127.8 (+, "},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"), 98.8 (C"},{"attributes":{"script":"sub"},"insert":"quat"},{"insert":", "},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"Ar"},{"insert":"I).\nIn accordance with literature: D. Chen, G. Shi, H. Jiang, Y. Zhang, Y. Zhang, "},{"attributes":{"italic":true},"insert":"Org. Lett."},{"insert":", "},{"attributes":{"bold":true},"insert":"2016"},{"insert":", "},{"attributes":{"italic":true},"insert":"18"},{"insert":", 2130–2133.\n"}]},
    expected: ' count: 10/12',
  },
  'REPO-C-42310-Alex': {
    formula: 'C11H12IN',
    content: {"ops":[{"attributes":{"script":"super"},"insert":"13"},{"insert":"C NMR (100 MHz, CDCl"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":", ppm) δ = 189.9 (1C, 1 × "},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":"), 154.3 (1C, 1 × "},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":"), 149.1 (1C,,1 × "},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":"), 137.6 (+,1C, 1 × "},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"ar"},{"insert":"), 131.6 (+, 1C, 1 × "},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"ar"},{"insert":"), 122.7 (+, 1C, 1 × "},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"ar"},{"insert":"), 90.8 (1C, 1 × "},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":"), 54.9 (1C, 1 × "},{"attributes":{"italic":true},"insert":"C"},{"attributes":{"script":"sub"},"insert":"q"},{"insert":"), 24.0 (+, 2C, 2 × "},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":"), 16.3 (+, 1C, 1 × "},{"attributes":{"italic":true},"insert":"C"},{"insert":"H"},{"attributes":{"script":"sub"},"insert":"3"},{"insert":").\n"}]},
    expected: '',
  },
  'REPO-42592': {
    formula: 'C16H12N2O8',
    content: { ops: [
      {insert: '13C NMR (125 MHz, CDCl3, ppm) δ = 164.7 (CO2CH3), 147.0 (2C, Cq), 137.8 (2C, Cq), 134.4 (2C, CH), 132.0 (2C, Cq), 131.0 (2C, CH), 126.2 (2C, CH), 53.1 (2C, CO2CH3).' }
    ] },
    expected: ' count: 15/16',
  },
  'NJ-dot': {
    formula: 'C81H50N2O8',
    content: { ops: [
      {insert: '13C NMR (126 MHz, CDCl3): δ = 167.0 (Cq, COOEt), 146.2 (Cq), 145.7 (Cq), 145.5 (Cq), 141.9 (Cq), 134.8 (+, CH), 134.6 (+, CH), 134.5 (+, CH), 134.5 (+, CH), 134.3 (+, CH), 134.2 (+, CH), 134.2 (+, CH), 129.7 (+, CH), 127.9 (+, CH), 127.9 (+, CH), 127.6 (+, CH), 127.5 (+, CH), 127.3 (+, CH), 127.2 (+, CH), 126.6 (+, CH), 126.6 (+, CH), 126.4 (+, CH), 119.8 (Cq), 119.5 (Cq), 118.0 (Cq), 117.3 (Cq), 61.5 (–, COCH2), 14.8 (+, CH2CH3) . Missing signals (52C) due to signal overlap.' }
    ] },
    expected: ' count: 80/81',
  },
  'NJ-dot2': {
    formula: 'C5H4',
    content: { ops: [
      {insert: '13C NMR : δ = 167.0 (CH), 14.8 (+, CH2CH3).Missing signals (2C) due to signal overlap.' }
    ] },
    expected: ' count: 4/5',
  },
  'NJ-dot3': {
    formula: 'C5H4',
    content: { ops: [
      {insert: '13C NMR : δ = 167.0 (CH), 14.8. Missing signals  due to signal overlap.' }
    ] },
    expected: ' count: 2/5',
  },

  'NJ-comma': {
    formula: 'C5H4',
    content: { ops: [
      {insert: '13C NMR : δ = 167.0 (CH), 14.8 (+, CH2CH3), Missing signals (2C) due to signal overlap.' }
    ] },
    expected: ' count: 4/5',
  }
};

export { nmrData13C };
