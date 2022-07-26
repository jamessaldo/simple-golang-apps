//
//  DetailQuestion.swift
//  Conversa
//
//  Created by zy on 26/07/22.
//

import SwiftUI

struct DetailQuestion: View {
    var body: some View {
        ZStack {
            LinearGradient(gradient: Gradient(colors: [Color(hex: "#FFFCF4"), Color(hex: "#FFF4D9"), Color(hex: "#FFF7E1")]), startPoint: .top, endPoint: .bottom)
                .edgesIgnoringSafeArea(.vertical)
            ScrollView {
                VStack(alignment: .leading) {
                    Text("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.")
                        .font(.headline)
                        .bold()
                        .foregroundColor(Color(hex: "#383838"))
                    VStack(alignment: .leading) {
                        HStack(spacing:16){
                            Image(systemName: "person.circle.fill")
                                .resizable()
                                .scaledToFit()
                                .frame(width: 32, height: 32)
                                .foregroundColor(.gray)
                                .padding(.top, 16)
                            VStack(alignment: .leading) {
                                Text("Juliet Van")
                                    .font(.caption)
                                    .bold()
                                    .foregroundColor(Color(hex: "#383838"))
                                    .padding(.top, 16)
                                    .multilineTextAlignment(.leading)
                                Text("6 hours ago")
                                    .font(.caption)
                                    .foregroundColor(Color(hex: "#383838"))
                                    .multilineTextAlignment(.leading)
                            }
                            Spacer()
                            Button(action: {
                            }) {
                                Image(systemName: "heart.circle")
                                    .foregroundColor(.gray)
                                    .padding(.top, 16)
                            }
                        }.padding(.horizontal, 24)
                        Text("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")
                            .font(.caption)
                            .foregroundColor(Color(hex: "#484848"))
                            .padding(.horizontal,24)
                            .padding(.top, 8)
                            .multilineTextAlignment(.leading)
                    }
                    .frame(
                        minWidth: 0,
                        maxWidth: 327,
                        alignment: .topLeading
                    )
                    .clipped()
                    .padding(.bottom)
                    .background(Color.white)
                    .cornerRadius(20)
                    .shadow(color: Color(hex: "#000000").opacity(0.05), radius: 3, x: 0, y: 3)
                    
                    Text("Answers (\(20))")
                        .font(.headline)
                        .bold()
                        .foregroundColor(Color(hex: "#383838"))
                        .padding(.top, 10)
                    
                    VStack(alignment: .leading) {
                        ForEach(0..<20){ index in
                            HStack(alignment: .top){
                                VStack(alignment: .center){
                                    Image(systemName: "person.circle.fill")
                                        .resizable()
                                        .scaledToFit()
                                        .frame(width: 32, height: 32)
                                        .foregroundColor(.gray)
                                        .padding(.top, 16)
                                    if index != 19{
                                        DottedLine()
                                            .stroke(style: StrokeStyle(lineWidth: 1, dash: [2]))
                                            .frame(maxWidth: 1, maxHeight: .infinity)
                                    }
                                }
                                
                                VStack(alignment: .leading){
                                    HStack{
                                        VStack(alignment: .leading) {
                                            Text("Juliet Van")
                                                .font(.caption)
                                                .bold()
                                                .foregroundColor(Color(hex: "#383838"))
                                                .padding(.top, 16)
                                                .multilineTextAlignment(.leading)
                                            Text("6 hours ago")
                                                .font(.caption)
                                                .foregroundColor(Color(hex: "#383838"))
                                                .multilineTextAlignment(.leading)
                                        }
                                        Spacer()
                                        Button(action: {
                                        }) {
                                            Image(systemName: "heart.circle")
                                                .foregroundColor(.gray)
                                                .padding(.top, 16)
                                        }
                                    }
                                    Text("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")
                                        .font(.caption)
                                        .foregroundColor(Color(hex: "#484848"))
                                        .padding(.top, 8)
                                        .multilineTextAlignment(.leading)
                                }
                            }.padding(.horizontal, 24)
                        }
                    }
                    .frame(
                        minWidth: 0,
                        maxWidth: 327,
                        alignment: .topLeading
                    )
                    .clipped()
                    .padding(.bottom)
                    .background(Color.white)
                    .cornerRadius(20)
                    .shadow(color: Color(hex: "#000000").opacity(0.05), radius: 3, x: 0, y: 3)
                }.frame(maxWidth: .infinity)
            }.padding(.horizontal, 24).padding(.top, 10)
        }
    }
}

struct DetailQuestion_Previews: PreviewProvider {
    static var previews: some View {
        DetailQuestion()
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
