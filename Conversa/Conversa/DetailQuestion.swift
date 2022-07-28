//
//  DetailQuestion.swift
//  Conversa
//
//  Created by zy on 26/07/22.
//

import SwiftUI

struct DetailQuestion: View {
    private var id: Int
    
    init(_ id: Int){
        self.id = id
        UITextView.appearance().backgroundColor = .clear
        UITableView.appearance().backgroundColor = .clear
    }
    
    @State private var contentText = ""
    @State var createComment: Bool = false
    @State var result = Post()
    @State var isCommentEmpty: Bool = false
    var body: some View {
        ZStack {
            LinearGradient(gradient: Gradient(colors: [Color(hex: "#FFFCF4"), Color(hex: "#FFF4D9"), Color(hex: "#FFF7E1")]), startPoint: .top, endPoint: .bottom)
                .edgesIgnoringSafeArea(.vertical)
            ScrollView {
                VStack(alignment: .leading) {
                    Text(result.title ?? "Fetch error!")
                        .font(.headline)
                        .bold()
                        .foregroundColor(Color(hex: "#383838"))
                        .padding(.top, 16)
                    VStack(alignment: .leading) {
                        HStack(spacing:16){
                            Image(systemName: "person.circle.fill")
                                .resizable()
                                .scaledToFit()
                                .frame(width: 32, height: 32)
                                .foregroundColor(.gray)
                                .padding(.top, 16)
                            VStack(alignment: .leading) {
                                Text(result.creator ?? "Fetch error!")
                                    .font(.caption)
                                    .bold()
                                    .foregroundColor(Color(hex: "#383838"))
                                    .padding(.top, 16)
                                    .multilineTextAlignment(.leading)
                                Text(formatter_toString.string(from:formatter_toDate.date(from: result.created_at ?? "") ?? Date()))
                                    .font(.caption)
                                    .foregroundColor(Color(hex: "#383838"))
                                    .multilineTextAlignment(.leading)
                            }
                        }.padding(.horizontal, 24)
                        Text(result.description ?? "Fetch error!")
                            .font(.caption)
                            .foregroundColor(Color(hex: "#484848"))
                            .padding(.horizontal,24)
                            .padding(.top, 8)
                            .multilineTextAlignment(.leading)
                        HStack (spacing: 16){
                            Button(action: {createComment = true}) {
                                Image(systemName: "bubble.left.circle")
                                    .resizable()
                                    .frame(width: 24, height: 24)
                                    .foregroundColor(.gray)
                            }
                            Button(action: {
                            }) {
                                Image(systemName: "heart.circle")
                                    .resizable()
                                    .frame(width: 24, height: 24)
                                    .foregroundColor(.gray)
                            }
                        }.padding(.horizontal, 24).padding(.top, 8)
                    }
                    .frame(
                        maxWidth: .infinity,
                        alignment: .topLeading
                    )
                    .clipped()
                    .padding(.bottom)
                    .background(Color.white)
                    .cornerRadius(20)
                    .shadow(color: Color(hex: "#000000").opacity(0.05), radius: 3, x: 0, y: 3)
                    
                    Text("Answers (\(result.comments?.count ?? 0))")
                        .font(.headline)
                        .bold()
                        .foregroundColor(Color(hex: "#383838"))
                        .padding(.top, 10)
                    if (result.comments?.count ?? 0) > 0 {
                        VStack(alignment: .leading) {
                            ForEach(Array(zip(result.comments?.indices ?? [].indices, result.comments ?? [])), id: \.0){ index, item in
                                HStack(alignment: .top){
                                    VStack(alignment: .center){
                                        Image(systemName: "person.circle.fill")
                                            .resizable()
                                            .scaledToFit()
                                            .frame(width: 32, height: 32)
                                            .foregroundColor(.gray)
                                            .padding(.top, 16)
                                        DottedLine()
                                            .stroke(style: StrokeStyle(lineWidth: 1, dash: [2]))
                                            .frame(maxWidth: 1, maxHeight: self.getDottedHeight(index: index))
                                    }
                                    
                                    VStack(alignment: .leading){
                                        HStack{
                                            VStack(alignment: .leading) {
                                                Text(item.creator ?? "Fetch error!")
                                                    .font(.caption)
                                                    .bold()
                                                    .foregroundColor(Color(hex: "#383838"))
                                                    .padding(.top, 16)
                                                    .multilineTextAlignment(.leading)
                                                Text(formatter_toString.string(from:formatter_toDate.date(from: item.created_at ?? "") ?? Date()))
                                                    .font(.caption)
                                                    .foregroundColor(Color(hex: "#383838"))
                                                    .multilineTextAlignment(.leading)
                                            }
                                            Spacer()
                                            Button(action: {
                                            }) {
                                                Image(systemName: "heart.circle")
                                                    .resizable()
                                                    .frame(width: 24, height: 24)
                                                    .foregroundColor(.gray)
                                                    .padding(.top, 8)
                                            }
                                        }
                                        Text(item.content ?? "")
                                            .font(.caption)
                                            .foregroundColor(Color(hex: "#484848"))
                                            .padding(.top, 8)
                                            .multilineTextAlignment(.leading)
                                    }
                                }.padding(.horizontal, 24)
                            }
                        }
                        .frame(
                            maxWidth: .infinity
                        )
                        .clipped()
                        .padding(.bottom)
                        .background(Color.white)
                        .cornerRadius(20)
                        .shadow(color: Color(hex: "#000000").opacity(0.05), radius: 3, x: 0, y: 3)
                    }
                }
                .sheet(isPresented: $createComment) {
                    CreateCommentModal(
                        postId: id,
                        isPresented: $createComment,
                        content: $contentText,
                        post: $result
                    )
                    .onDisappear(perform: {
                        loadData()
                    })
                }
                .padding(.horizontal, 30).padding(.top, 10)
            }
        }.onAppear(perform: loadData)
    }
    
    private func loadData() {
        guard let url = URL(string: "http://127.0.0.1:8888/post/\(self.id)") else {
            print("Invalid URL")
            return
        }
        let request = URLRequest(url: url)
        
        URLSession.shared.dataTask(with: request) { data, response, error in
            if let data = data {
                if let response = try? JSONDecoder().decode(Post.self, from: data) {
                    DispatchQueue.main.async {
                        self.result = response
                        self.isCommentEmpty = response.comments?.count ?? 0 == 0
                    }
                    return
                }
            }
        }.resume()
    }
    
    private func getDottedHeight(index: Int) -> CGFloat {
        if index != (result.comments?.count ?? 0) - 1 {
            return .infinity
        }
        return 0
    }
}

struct DetailQuestion_Previews: PreviewProvider {
    static var previews: some View {
        DetailQuestion(1)
    }
}

struct DottedLine: Shape {
    
    func path(in rect: CGRect) -> Path {
        var path = Path()
        path.move(to: CGPoint(x: 0, y: 0))
        path.addLine(to: CGPoint(x: rect.width, y: rect.height))
        return path
    }
}
