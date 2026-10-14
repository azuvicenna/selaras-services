#set page(
  paper: "a4",
  margin: (x: 2cm, y: 2.2cm),
  header: [
    #align(center)[
      #text(13pt, weight: "bold")[RUMAH SAKIT SELARAS MEDIKA] \
      #text(8pt)[Jl. Kesehatan Sejahtera No. 100, Jakarta | Telp: (021) 555-0123]
      #v(2pt)
      #line(length: 100%, stroke: 1.5pt)
    ]
  ],
  footer: [
    #line(length: 100%, stroke: 0.5pt)
    #grid(
      columns: (1fr, 1fr),
      align(left)[#text(7pt, fill: luma(100))[Formulir Penolakan Tindakan Kedokteran (Informed Refusal)]],
      align(right)[#text(7pt)[Halaman #context counter(page).display()]]
    )
  ]
)

#set text(size: 9.5pt)

#align(center)[
  #text(12pt, weight: "bold")[SURAT PENOLAKAN TINDAKAN MEDIS (INFORMED REFUSAL)] \
  #text(8.5pt)[Nomor: {{DOC_NUMBER}}]
]

#v(0.8em)

Yang bertanda tangan di bawah ini:
#table(
  columns: (130pt, 8pt, 1fr),
  stroke: none,
  [Nama], [:], [{{PENOLAK_NAME}}],
  [Jenis Kelamin / Tgl Lahir], [:], [{{PENOLAK_GENDER}} / {{PENOLAK_DOB}}],
  [Alamat Domisili], [:], [{{PENOLAK_ALAMAT}}],
  [Hubungan dengan Pasien], [:], [{{HUBUNGAN}}],
)

Dengan ini menyatakan *MENOLAK* untuk dilakukannya tindakan kedokteran berupa:
#align(center)[
  #rect(stroke: 1pt, inset: 8pt)[
    #text(11pt, weight: "bold")[{{TINDAKAN_DITOLAK}}]
  ]
]

Terhadap pasien dengan identitas:
#table(
  columns: (130pt, 8pt, 1fr),
  stroke: none,
  [Nama Pasien], [:], [{{PATIENT_NAME}}],
  [Nomor Rekam Medis], [:], [{{PATIENT_ID}}],
  [Diagnosa Dokter], [:], [{{DIAGNOSIS}}],
)

#v(0.5em)
Saya menyatakan bahwa saya telah menerima penjelasan lengkap dari dokter mengenai:
1. Diagnosis dan alasan diperlukannya tindakan medis tersebut.
2. *Risiko dan akibat penolakan*, termasuk kemungkinan terjadinya perburukan kondisi kesehatan, cacat permanen, hingga *risiko kematian*.
3. Konsekuensi hukum dan biaya yang timbul atas keputusan penolakan ini menjadi tanggung jawab pribadi saya/keluarga sepenuhnya, serta membebaskan pihak Rumah Sakit Selaras Medika dan tenaga medis dari tuntutan hukum di kemudian hari.

#v(1.5em)

#grid(
  columns: (1fr, 1fr, 1fr),
  [
    #align(center)[
      Dokter Pemeriksa, \
      #v(3.5em)
      #text(weight: "bold")[{{DOCTOR_NAME}}]
    ]
  ],
  [
    #align(center)[
      Saksi Pihak RS, \
      #v(3.5em)
      #text(weight: "bold")[{{SAKSI_RS}}]
    ]
  ],
  [
    #align(center)[
      Yang Menyatakan Menolak, \
      #v(3.5em)
      #text(weight: "bold")[{{PENOLAK_NAME}}]
    ]
  ]
)
