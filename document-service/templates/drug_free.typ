#set page(
  paper: "a4",
  margin: (x: 2cm, y: 2.2cm),
  header: [
    #align(center)[
      #text(13pt, weight: "bold")[RUMAH SAKIT SELARAS MEDIKA] \
      #text(10pt, weight: "bold")[TIM PEMERIKSA KESEHATAN & BEBAS NARKOBA] \
      #text(8pt)[Jl. Kesehatan Sejahtera No. 100, Jakarta | Telp: (021) 555-0123]
      #v(2pt)
      #line(length: 100%, stroke: 1.5pt)
    ]
  ],
  footer: [
    #line(length: 100%, stroke: 0.5pt)
    #grid(
      columns: (1fr, 1fr),
      align(left)[#text(7pt, fill: luma(100))[Surat Keterangan Bebas Narkoba (SKBN) Resmi - RS Selaras Medika]],
      align(right)[#text(7pt)[Halaman #context counter(page).display()]]
    )
  ]
)

#set text(size: 9.5pt)

#align(center)[
  #text(12pt, weight: "bold")[SURAT KETERANGAN BEBAS NARKOBA (SKBN)] \
  #text(9pt)[Nomor: {{DOC_NUMBER}}]
]

#v(0.8em)

Yang bertanda tangan di bawah ini Dokter Pemeriksa Rumah Sakit Selaras Medika menerangkan bahwa:

#table(
  columns: (140pt, 10pt, 1fr),
  stroke: none,
  [Nama Lengkap], [:], [{{PATIENT_NAME}}],
  [Nomor Rekam Medis], [:], [{{PATIENT_ID}}],
  [Nomor KTP / NIK], [:], [{{NIK}}],
  [Tempat, Tanggal Lahir], [:], [{{PATIENT_POB}}, {{PATIENT_DOB}}],
  [Jenis Kelamin / Pekerjaan], [:], [{{PATIENT_GENDER}} / {{PEKERJAAN}}],
  [Alamat Domisili], [:], [{{ALAMAT}}],
)

#v(0.5em)
Telah dilakukan pemeriksaan fisik dan skrining laboratorium toksikologi urine (*Multi-Panel Drug Screening*) pada tanggal *{{ISSUED_DATE}}* dengan hasil sebagai berikut:

#table(
  columns: (30pt, 160pt, 100pt, 1fr),
  stroke: 0.5pt,
  table.header([*No*], [*Parameter Zat Uji (NAPZA)*], [*Hasil Uji*], [*Keterangan*]),
  [1.], [Amphetamine (AMP)], [{{AMP_RESULT}}], [Dalam batas cut-off normal],
  [2.], [Methamphetamine / Shabu (MET)], [{{MET_RESULT}}], [Dalam batas cut-off normal],
  [3.], [THC / Ganja / Marijuana (THC)], [{{THC_RESULT}}], [Dalam batas cut-off normal],
  [4.], [Morphine / Opiates (MOP)], [{{MOP_RESULT}}], [Dalam batas cut-off normal],
  [5.], [Benzodiazepine (BZO)], [{{BZO_RESULT}}], [Dalam batas cut-off normal],
  [6.], [Cocaine (COC)], [{{COC_RESULT}}], [Dalam batas cut-off normal],
)

#v(0.8em)
*KESIMPULAN:*
Berdasarkan hasil pemeriksaan tersebut di atas, pada saat diperiksa dinyatakan:
#align(center)[
  #rect(stroke: 1pt, inset: 8pt)[
    #text(11pt, weight: "bold")[BEBAS DARI ZAT NARKOTIKA & PSIKOTROPIKA (NAPZA)]
  ]
]

Surat keterangan ini diterbitkan untuk keperluan: *{{KEPERLUAN}}*.

#v(1.5em)

#grid(
  columns: (1fr, 1fr),
  [
    #rect(stroke: 0.5pt, inset: 6pt)[
      #align(center)[
        #text(7pt, weight: "bold")[VERIFIKASI KEASLIAN ELEKTRONIK] \
        #text(6.5pt)[Pindai kode QR untuk validasi keaslian dokumen:] \
        #v(2pt)
        #text(6.5pt, font: ("Courier New", "DejaVu Sans Mono", "monospace"))[{{VERIFY_URL}}]
      ]
    ]
  ],
  [
    #align(center)[
      Jakarta, {{ISSUED_DATE}} \
      Dokter Pemeriksa, \
      #v(3.5em)
      #text(weight: "bold")[{{DOCTOR_NAME}}] \
      SIP: {{DOCTOR_SIP}}
    ]
  ]
)
