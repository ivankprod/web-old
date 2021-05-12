/* 
    UTILS SCRIPT

    Author : IvanK Production
*/


//  Get meta tag
function getMeta(doc, metaName) {
	let result = doc.querySelector('meta[name="' + metaName + '"]');
	if (!result) result = doc.querySelector('meta[property="' + metaName + '"]');
	if (!result) throw new Error("Meta not found!");

	return result.content;
}

//  Set meta tag
function setMeta(doc, metaName, metaContent) {
	let result = doc.querySelector('meta[name="' + metaName + '"]');
	if (!result) result = doc.querySelector('meta[property="' + metaName + '"]');
	if (!result) throw new Error("Meta not found!");

	result.content = metaContent;
}

//  Rewrite meta tags from ajax-loaded page to current page
function rewriteMetas(opts) {
	if (opts && opts.metas && opts.docSource && opts.docDest) {
		opts.metas.forEach(elem => {
			setMeta(opts.docDest, elem, getMeta(opts.docSource, elem));
		})
	} else {
		throw new Error("Options not specified!");
	}
}

//  Exports
export default {
	rewriteMetas
}
