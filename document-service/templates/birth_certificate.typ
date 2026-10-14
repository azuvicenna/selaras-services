#set page(
  paper: "a4",
  margin: (x: 2cm, y: 2.2cm),
  header: [
    #align(center)[
      #text(13pt, weight: "bold")[RUMAH SAKIT SELARAS MEDIKA] \
      #text(10pt, weight: "bold")[INSTALASI KEBIDANAN & KANDUNGAN (VK)] \
      #text(8pt)[Jl. Kesehatan Sejahtera No. 100, Jakarta | Telp: (021) 555-0123]
      #v(2pt)
      #line(length: 100%, stroke: 1.5pt)
    ]
  ],
  footer: [
    #line(length: 100%, stroke: 0.5pt)
    #grid(
      columns: (1fr, 1fr),
      align(left)[#text(7pt, fill: luma(100))[Surat Keterangan Kelahiran Resmi - RS Selaras Medika]],
      align(right)[#text(7pt)[Halaman #context counter(page).display()]]
    )
  ]
)

#set text(size: 9.5pt)

#align(center)[
  #text(12pt, weight: "bold")[SURAT KETERANGAN KELAHIRAN] \
  #text(9pt)[Nomor: {{DOC_NUMBER}}]
]

#v(0.8em)

Yang bertanda tangan di bawah ini menerangkan bahwa telah lahir seorang bayi pada:

#table(
  columns: (140pt, 10pt, 1fr),
  stroke: none,
  [Hari / Tanggal Lahir], [:], [{{HARI_LAHIR}}, {{TGL_LAHIR}}],
  [Jam / Waktu Kelahiran], [:], [{{JAM_LAHIR}} WIB],
  [Tempat Kelahiran], [:], [Kamar Bersalin RS Selaras Medika Jakarta],
)

#v(0.5em)
*DATA BAYI:*
#table(
  columns: (140pt, 10pt, 1fr),
  stroke: 0.5pt,
  [Nama Bayi], [:], [{{NAMA_BAYI}}],
  [Jenis Kelamin], [:], [{{JENIS_KELAMIN_BAYI}}],
  [Berat Badan Lahir], [:], [{{BERAT_BADAN}} gram],
  [Panjang Badan Lahir], [:], [{{PANJANG_BADAN}} cm],
  [Lingkar Kepala / Dada], [:], [{{LINGKAR_KEPALA}} cm / {{LINGKAR_DADA}} cm],
  [Nilai APGAR], [:], [{{APGAR_SCORE}}],
  [Kelahiran ke-], [:], [{{KELAHIRAN_KE}} (Tunggal / Gemelli)],
)

#v(0.5em)
*DATA ORANG TUA:*
#table(
  columns: (140pt, 10pt, 1fr),
  stroke: 0.5pt,
  [Nama Ibu Kandung], [:], [{{PATIENT_NAME}} (NIK: {{NIK_IBU}})],
  [Nomor RM Ibu], [:], [{{PATIENT_ID}}],
  [Nama Ayah Kandung], [:], [{{NAMA_AYAH}} (NIK: {{NIK_AYAH}})],
  [Alamat Domisili], [:], [{{ALAMAT}}],
)

#v(1em)
Surat keterangan kelahiran ini dibuat sebagai bukti otentik kelahiran untuk pengurusan Akta Kelahiran dan dokumen kependudukan resmi.

#v(1.5em)

#grid(
  columns: (1fr, 1fr),
  [
    #rect(stroke: 0.5pt, inset: 6pt)[
      #align(center)[
        #text(7pt, weight: "bold")[VERIFIKASI KEASLIAN ELEKTRONIK] \
        #text(6.5pt)[Pindai kode QR untuk validasi data kelahiran:] \
        #v(2pt)
        #text(6.5pt, font: ("Courier New", "DejaVu Sans Mono", "monospace"))[{{VERIFY_URL}}]
      ]
    ]
  ],
  [
    #align(center)[
      Jakarta, {{ISSUED_DATE}} \
      Dokter Penolong Persalinan (Sp.OG), \
      #v(3.5em)
      #text(weight: "bold")[{{DOCTOR_NAME}}] \
      SIP: {{DOCTOR_SIP}}
    ]
  ]
)
