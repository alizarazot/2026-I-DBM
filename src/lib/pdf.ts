import jsPDF from "jspdf";
import autoTable from "jspdf-autotable";

export function downloadPdf(
	title: string,
	headings: string[],
	rows: (string | number)[][],
) {
	const doc = new jsPDF();
	doc.setFontSize(18);
	doc.text(title, 14, 22);
	autoTable(doc, {
		head: [headings],
		body: rows,
		startY: 30,
		styles: { fontSize: 10 },
	});
	doc.save(`${title.toLowerCase().replace(/\s+/g, "-")}.pdf`);
}
